#include <napi.h>
#include <open62541/client.h>
#include <open62541/client_config_default.h>
#include <open62541/client_highlevel.h>
#include <open62541/client_subscriptions.h>

#include <atomic>
#include <chrono>
#include <map>
#include <memory>
#include <mutex>
#include <string>
#include <thread>
#include <vector>

// C++20 N-API wrapper for open62541 client operations.
// The N-API boundary uses modern C++ (C++20); all open62541 calls remain C-style.

namespace {

// ---------------------------------------------------------------------------
// Helpers: status codes, NodeId parsing, data type conversion
// ---------------------------------------------------------------------------

static std::string StatusCodeString(UA_StatusCode code) {
  if (UA_StatusCode_isGood(code)) return "Good";
  const char* name = UA_StatusCode_name(code);
  return name ? std::string(name) : "Unknown";
}

static Napi::Value DateTimeToJs(Napi::Env env, UA_DateTime dt) {
  // UA_DateTime is 100 ns intervals since 1601-01-01 UTC.
  // Convert to Unix milliseconds.
  static constexpr UA_Int64 epochOffset = 116444736000000000LL;  // 100 ns
  UA_Int64 unixHns = dt - epochOffset;
  double unixMs = static_cast<double>(unixHns) / 10000.0;
  return Napi::Date::New(env, unixMs);
}

static bool ParseNodeId(const std::string& str, UA_NodeId* out) {
  UA_NodeId_init(out);
  UA_String s = UA_STRING((char*)str.c_str());
  UA_StatusCode status = UA_NodeId_parse(out, s);
  return UA_StatusCode_isGood(status);
}

static const UA_DataType* DataTypeFromName(const std::string& name) {
  if (name == "Boolean") return &UA_TYPES[UA_TYPES_BOOLEAN];
  if (name == "SByte") return &UA_TYPES[UA_TYPES_SBYTE];
  if (name == "Byte") return &UA_TYPES[UA_TYPES_BYTE];
  if (name == "Int16") return &UA_TYPES[UA_TYPES_INT16];
  if (name == "UInt16") return &UA_TYPES[UA_TYPES_UINT16];
  if (name == "Int32") return &UA_TYPES[UA_TYPES_INT32];
  if (name == "UInt32") return &UA_TYPES[UA_TYPES_UINT32];
  if (name == "Int64") return &UA_TYPES[UA_TYPES_INT64];
  if (name == "UInt64") return &UA_TYPES[UA_TYPES_UINT64];
  if (name == "Float") return &UA_TYPES[UA_TYPES_FLOAT];
  if (name == "Double") return &UA_TYPES[UA_TYPES_DOUBLE];
  if (name == "String") return &UA_TYPES[UA_TYPES_STRING];
  if (name == "DateTime") return &UA_TYPES[UA_TYPES_DATETIME];
  if (name == "ByteString") return &UA_TYPES[UA_TYPES_BYTESTRING];
  if (name == "LocalizedText") return &UA_TYPES[UA_TYPES_LOCALIZEDTEXT];
  if (name == "QualifiedName") return &UA_TYPES[UA_TYPES_QUALIFIEDNAME];
  return nullptr;
}

static std::string DataTypeName(const UA_DataType* type) {
  if (!type) return "Unknown";
  if (type == &UA_TYPES[UA_TYPES_BOOLEAN]) return "Boolean";
  if (type == &UA_TYPES[UA_TYPES_SBYTE]) return "SByte";
  if (type == &UA_TYPES[UA_TYPES_BYTE]) return "Byte";
  if (type == &UA_TYPES[UA_TYPES_INT16]) return "Int16";
  if (type == &UA_TYPES[UA_TYPES_UINT16]) return "UInt16";
  if (type == &UA_TYPES[UA_TYPES_INT32]) return "Int32";
  if (type == &UA_TYPES[UA_TYPES_UINT32]) return "UInt32";
  if (type == &UA_TYPES[UA_TYPES_INT64]) return "Int64";
  if (type == &UA_TYPES[UA_TYPES_UINT64]) return "UInt64";
  if (type == &UA_TYPES[UA_TYPES_FLOAT]) return "Float";
  if (type == &UA_TYPES[UA_TYPES_DOUBLE]) return "Double";
  if (type == &UA_TYPES[UA_TYPES_STRING]) return "String";
  if (type == &UA_TYPES[UA_TYPES_DATETIME]) return "DateTime";
  if (type == &UA_TYPES[UA_TYPES_BYTESTRING]) return "ByteString";
  if (type == &UA_TYPES[UA_TYPES_LOCALIZEDTEXT]) return "LocalizedText";
  if (type == &UA_TYPES[UA_TYPES_QUALIFIEDNAME]) return "QualifiedName";
  return type->typeName ? std::string(type->typeName) : "Unknown";
}

static bool JsValueToVariantScalar(Napi::Env env, const Napi::Value& value,
                                   const std::string& dataTypeName,
                                   UA_Variant* outVariant) {
  UA_Variant_init(outVariant);
  const UA_DataType* type = DataTypeFromName(dataTypeName);
  if (!type) {
    Napi::Error::New(env, "Unsupported data type: " + dataTypeName).ThrowAsJavaScriptException();
    return false;
  }

  if (dataTypeName == "Boolean") {
    UA_Boolean v = value.As<Napi::Boolean>().Value();
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "SByte") {
    UA_SByte v = static_cast<UA_SByte>(value.As<Napi::Number>().Int32Value());
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "Byte") {
    UA_Byte v = static_cast<UA_Byte>(value.As<Napi::Number>().Uint32Value() & 0xFF);
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "Int16") {
    UA_Int16 v = static_cast<UA_Int16>(value.As<Napi::Number>().Int32Value());
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "UInt16") {
    UA_UInt16 v = static_cast<UA_UInt16>(value.As<Napi::Number>().Uint32Value() & 0xFFFF);
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "Int32") {
    UA_Int32 v = value.As<Napi::Number>().Int32Value();
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "UInt32") {
    UA_UInt32 v = value.As<Napi::Number>().Uint32Value();
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "Int64") {
    UA_Int64 v = static_cast<UA_Int64>(value.As<Napi::Number>().Int64Value());
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "UInt64") {
    UA_UInt64 v = static_cast<UA_UInt64>(value.As<Napi::Number>().Int64Value());
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "Float") {
    UA_Float v = value.As<Napi::Number>().FloatValue();
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "Double") {
    UA_Double v = value.As<Napi::Number>().DoubleValue();
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "String") {
    std::string s = value.As<Napi::String>().Utf8Value();
    UA_String us = UA_STRING((char*)s.c_str());
    UA_Variant_setScalarCopy(outVariant, &us, type);
  } else if (dataTypeName == "DateTime") {
    double ms;
    if (value.IsDate()) {
      ms = value.As<Napi::Date>().ValueOf();
    } else {
      // For ISO strings, fall back to current time in this MVP.
      ms = static_cast<double>(
          std::chrono::duration_cast<std::chrono::milliseconds>(
              std::chrono::system_clock::now().time_since_epoch())
              .count());
    }
    static constexpr UA_Int64 epochOffset = 116444736000000000LL;
    UA_DateTime v = static_cast<UA_DateTime>(ms * 10000.0) + epochOffset;
    UA_Variant_setScalarCopy(outVariant, &v, type);
  } else if (dataTypeName == "ByteString") {
    std::string s = value.As<Napi::String>().Utf8Value();
    UA_ByteString bs;
    UA_ByteString_init(&bs);
    bs.length = s.size();
    bs.data = (UA_Byte*)UA_malloc(bs.length);
    if (bs.data && bs.length) memcpy(bs.data, s.data(), bs.length);
    UA_Variant_setScalarCopy(outVariant, &bs, type);
    UA_ByteString_clear(&bs);
  } else if (dataTypeName == "LocalizedText") {
    std::string s = value.As<Napi::String>().Utf8Value();
    UA_LocalizedText lt = UA_LOCALIZEDTEXT((char*)"en-US", (char*)s.c_str());
    UA_Variant_setScalarCopy(outVariant, &lt, type);
  } else if (dataTypeName == "QualifiedName") {
    std::string s = value.As<Napi::String>().Utf8Value();
    UA_QualifiedName qn = UA_QUALIFIEDNAME(1, (char*)s.c_str());
    UA_Variant_setScalarCopy(outVariant, &qn, type);
  } else {
    Napi::Error::New(env, "Unsupported data type: " + dataTypeName).ThrowAsJavaScriptException();
    return false;
  }
  return true;
}

static Napi::Value VariantToJsValue(Napi::Env env, const UA_Variant* variant) {
  if (!variant || !variant->type || !variant->data) {
    return env.Null();
  }

  const UA_DataType* type = variant->type;
  void* data = variant->data;
  bool isArray = variant->arrayLength > 0;

  auto scalarValue = [&]() -> Napi::Value {
    if (type == &UA_TYPES[UA_TYPES_BOOLEAN]) {
      return Napi::Boolean::New(env, *(UA_Boolean*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_SBYTE]) {
      return Napi::Number::New(env, *(UA_SByte*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_BYTE]) {
      return Napi::Number::New(env, *(UA_Byte*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_INT16]) {
      return Napi::Number::New(env, *(UA_Int16*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_UINT16]) {
      return Napi::Number::New(env, *(UA_UInt16*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_INT32]) {
      return Napi::Number::New(env, *(UA_Int32*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_UINT32]) {
      return Napi::Number::New(env, *(UA_UInt32*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_INT64]) {
      return Napi::Number::New(env, static_cast<double>(*(UA_Int64*)data));
    }
    if (type == &UA_TYPES[UA_TYPES_UINT64]) {
      return Napi::Number::New(env, static_cast<double>(*(UA_UInt64*)data));
    }
    if (type == &UA_TYPES[UA_TYPES_FLOAT]) {
      return Napi::Number::New(env, *(UA_Float*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_DOUBLE]) {
      return Napi::Number::New(env, *(UA_Double*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_STRING]) {
      UA_String* s = (UA_String*)data;
      return Napi::String::New(env, s->data ? std::string((char*)s->data, s->length) : "");
    }
    if (type == &UA_TYPES[UA_TYPES_DATETIME]) {
      return DateTimeToJs(env, *(UA_DateTime*)data);
    }
    if (type == &UA_TYPES[UA_TYPES_BYTESTRING]) {
      UA_ByteString* bs = (UA_ByteString*)data;
      return Napi::String::New(env,
                               bs->data ? std::string((char*)bs->data, bs->length) : "");
    }
    if (type == &UA_TYPES[UA_TYPES_LOCALIZEDTEXT]) {
      UA_LocalizedText* lt = (UA_LocalizedText*)data;
      Napi::Object obj = Napi::Object::New(env);
      obj.Set("locale",
              lt->locale.data ? std::string((char*)lt->locale.data, lt->locale.length) : "");
      obj.Set("text",
              lt->text.data ? std::string((char*)lt->text.data, lt->text.length) : "");
      return obj;
    }
    if (type == &UA_TYPES[UA_TYPES_QUALIFIEDNAME]) {
      UA_QualifiedName* qn = (UA_QualifiedName*)data;
      Napi::Object obj = Napi::Object::New(env);
      obj.Set("namespaceIndex", qn->namespaceIndex);
      obj.Set("name",
              qn->name.data ? std::string((char*)qn->name.data, qn->name.length) : "");
      return obj;
    }
    return env.Null();
  };

  if (!isArray) {
    return scalarValue();
  }

  Napi::Array arr = Napi::Array::New(env, variant->arrayLength);
  size_t typeSize = type->memSize;
  for (size_t i = 0; i < variant->arrayLength; i++) {
    UA_Variant item;
    UA_Variant_init(&item);
    item.type = type;
    item.storageType = UA_VARIANT_DATA_NODELETE;
    item.data = (char*)data + i * typeSize;
    item.arrayLength = 0;
    arr.Set(i, VariantToJsValue(env, &item));
  }
  return arr;
}

static std::string SecurityPolicyUri(const std::string& name) {
  if (name == "None") return "http://opcfoundation.org/UA/SecurityPolicy#None";
  if (name == "Basic128Rsa15") return "http://opcfoundation.org/UA/SecurityPolicy#Basic128Rsa15";
  if (name == "Basic256") return "http://opcfoundation.org/UA/SecurityPolicy#Basic256";
  if (name == "Basic256Sha256") return "http://opcfoundation.org/UA/SecurityPolicy#Basic256Sha256";
  if (name == "Aes128Sha256RsaOaep")
    return "http://opcfoundation.org/UA/SecurityPolicy#Aes128Sha256RsaOaep";
  if (name == "Aes256Sha256RsaPss")
    return "http://opcfoundation.org/UA/SecurityPolicy#Aes256Sha256RsaPss";
  return "";
}

static UA_MessageSecurityMode SecurityModeFromString(const std::string& mode) {
  if (mode == "Sign") return UA_MESSAGESECURITYMODE_SIGN;
  if (mode == "SignAndEncrypt") return UA_MESSAGESECURITYMODE_SIGNANDENCRYPT;
  return UA_MESSAGESECURITYMODE_NONE;
}

// ---------------------------------------------------------------------------
// Client context and management
// ---------------------------------------------------------------------------

struct SubscriptionContext {
  Napi::ThreadSafeFunction tsfn;
  std::string nodeId;
  uint32_t subscriptionId = 0;
  uint32_t monitoredItemId = 0;
};

struct ClientContext {
  UA_Client* client = nullptr;
  uint64_t id = 0;
  std::atomic<bool> running{false};
  std::thread iterateThread;
  std::mutex clientMutex;
  std::map<uint32_t, std::unique_ptr<SubscriptionContext>> subscriptions;
};

static std::mutex g_clients_mutex;
static std::map<uint64_t, std::shared_ptr<ClientContext>> g_clients;
static std::atomic<uint64_t> g_next_client_id{1};

static void ClientIterateLoop(ClientContext* ctx) {
  while (ctx->running.load()) {
    std::lock_guard<std::mutex> lock(ctx->clientMutex);
    if (!ctx->client) break;
    UA_Client_run_iterate(ctx->client, 100);
  }
}

static void PauseClientIteration(const std::shared_ptr<ClientContext>& ctx) {
  ctx->running.store(false);
  if (ctx->iterateThread.joinable()) {
    ctx->iterateThread.join();
  }
}

static void ResumeClientIteration(const std::shared_ptr<ClientContext>& ctx) {
  ctx->running.store(true);
  ctx->iterateThread = std::thread(ClientIterateLoop, ctx.get());
}

static void DeleteClientContext(uint64_t clientId) {
  std::shared_ptr<ClientContext> ctx;
  {
    std::lock_guard<std::mutex> lock(g_clients_mutex);
    auto it = g_clients.find(clientId);
    if (it == g_clients.end()) return;
    ctx = std::move(it->second);
    g_clients.erase(it);
  }

  if (!ctx) return;

  ctx->running.store(false);
  if (ctx->iterateThread.joinable()) {
    ctx->iterateThread.join();
  }

  // Release TSFNs before deleting the client.
  for (auto& pair : ctx->subscriptions) {
    if (pair.second && pair.second->tsfn) {
      pair.second->tsfn.Release();
    }
  }
  ctx->subscriptions.clear();

  if (ctx->client) {
    UA_Client_disconnect(ctx->client);
    UA_Client_delete(ctx->client);
    ctx->client = nullptr;
  }
}

static void DataChangeCallback(UA_Client* client, UA_UInt32 subId, void* subContext,
                               UA_UInt32 monId, void* monContext, UA_DataValue* value) {
  (void)client;
  (void)subId;
  (void)subContext;
  (void)monId;

  SubscriptionContext* sub = static_cast<SubscriptionContext*>(monContext);
  if (!sub || !sub->tsfn) return;

  struct CallbackData {
    std::string nodeId;
    std::string dataType;
    std::string statusCode;
    double sourceTimestamp = 0;
    double serverTimestamp = 0;
    bool hasSourceTimestamp = false;
    bool hasServerTimestamp = false;
    Napi::Value value;  // marshalled on the JS thread
  };

  // We cannot create Napi::Value here. Marshal the raw variant into a heap copy
  // and convert it in the TSFN callback on the JS thread.
  struct RawCallbackData {
    std::string nodeId;
    UA_Variant variant;
    UA_DataValue dataValue;
  };

  RawCallbackData* raw = new RawCallbackData();
  raw->nodeId = sub->nodeId;
  UA_Variant_init(&raw->variant);
  UA_DataValue_init(&raw->dataValue);

  if (value && value->hasValue && value->value.type && value->value.data) {
    UA_Variant_copy(&value->value, &raw->variant);
    raw->dataValue.hasSourceTimestamp = value->hasSourceTimestamp;
    raw->dataValue.sourceTimestamp = value->sourceTimestamp;
    raw->dataValue.hasServerTimestamp = value->hasServerTimestamp;
    raw->dataValue.serverTimestamp = value->serverTimestamp;
  }

  sub->tsfn.NonBlockingCall(raw, [](Napi::Env env, Napi::Function jsCallback,
                                   RawCallbackData* rawData) {
    Napi::Object msg = Napi::Object::New(env);
    msg.Set("nodeId", rawData->nodeId);
    msg.Set("value", VariantToJsValue(env, &rawData->variant));
    msg.Set("dataType", DataTypeName(rawData->variant.type));
    msg.Set("statusCode", StatusCodeString(UA_STATUSCODE_GOOD));

    if (rawData->dataValue.hasSourceTimestamp) {
      msg.Set("sourceTimestamp", DateTimeToJs(env, rawData->dataValue.sourceTimestamp));
    }
    if (rawData->dataValue.hasServerTimestamp) {
      msg.Set("serverTimestamp", DateTimeToJs(env, rawData->dataValue.serverTimestamp));
    }

    jsCallback.Call({env.Null(), msg});
    UA_Variant_clear(&rawData->variant);
    delete rawData;
  });
}

// ---------------------------------------------------------------------------
// Connection options parsing
// ---------------------------------------------------------------------------

struct ConnectionOptions {
  std::string endpoint;
  std::string securityMode;
  std::string securityPolicy;
  std::string user;
  std::string password;
  std::string certificate;
  std::string privateKey;
  int requestedSessionTimeout = 60000;
};

static std::string GetStringOrDefault(const Napi::Object& opts, const char* key,
                                      const std::string& defaultValue = "") {
  if (!opts.Has(key)) return defaultValue;
  Napi::Value value = opts.Get(key);
  if (value.IsUndefined() || value.IsNull() || !value.IsString()) return defaultValue;
  return value.As<Napi::String>().Utf8Value();
}

static int GetIntOrDefault(const Napi::Object& opts, const char* key, int defaultValue) {
  if (!opts.Has(key)) return defaultValue;
  Napi::Value value = opts.Get(key);
  if (value.IsUndefined() || value.IsNull() || !value.IsNumber()) return defaultValue;
  return value.As<Napi::Number>().Int32Value();
}

static ConnectionOptions ParseConnectionOptions(const Napi::Object& opts) {
  ConnectionOptions result;
  result.endpoint = GetStringOrDefault(opts, "endpoint");
  result.securityMode = GetStringOrDefault(opts, "securityMode");
  result.securityPolicy = GetStringOrDefault(opts, "securityPolicy");
  result.user = GetStringOrDefault(opts, "user");
  result.password = GetStringOrDefault(opts, "password");
  result.certificate = GetStringOrDefault(opts, "certificate");
  result.privateKey = GetStringOrDefault(opts, "privateKey");
  result.requestedSessionTimeout = GetIntOrDefault(opts, "requestedSessionTimeout", 60000);
  return result;
}

// ---------------------------------------------------------------------------
// AsyncWorkers
// ---------------------------------------------------------------------------

class ConnectWorker : public Napi::AsyncWorker {
 public:
  ConnectWorker(Napi::Promise::Deferred deferred, const ConnectionOptions& options)
      : Napi::AsyncWorker(deferred.Env()),
        deferred_(deferred),
        options_(options),
        client_id_(0),
        status_(UA_STATUSCODE_GOOD) {}

  void Execute() override {
    std::shared_ptr<ClientContext> ctx = std::make_shared<ClientContext>();
    ctx->client = UA_Client_new();
    if (!ctx->client) {
      SetError("Failed to create open62541 client");
      return;
    }

    UA_ClientConfig* config = UA_Client_getConfig(ctx->client);
    UA_ClientConfig_setDefault(config);
    config->securityMode = SecurityModeFromString(options_.securityMode);
    std::string policyUri = SecurityPolicyUri(options_.securityPolicy);
    if (!policyUri.empty()) {
      UA_String_clear(&config->securityPolicyUri);
      config->securityPolicyUri = UA_STRING_ALLOC(policyUri.c_str());
    }
    config->timeout = 10000;

    // TODO: certificate loading for Sign/SignAndEncrypt modes is left as a
    // follow-up. For None mode the client connects without certificates.
    if (!options_.certificate.empty() || !options_.privateKey.empty()) {
      // Best-effort placeholder: real cert loading requires UA_ClientConfig_setDefaultEncryption
      // and PKI configuration.
    }

    UA_StatusCode status;
    if (!options_.user.empty()) {
      status = UA_Client_connectUsername(ctx->client, options_.endpoint.c_str(),
                                         options_.user.c_str(), options_.password.c_str());
    } else {
      status = UA_Client_connect(ctx->client, options_.endpoint.c_str());
    }

    if (!UA_StatusCode_isGood(status)) {
      status_ = status;
      UA_Client_delete(ctx->client);
      ctx->client = nullptr;
      SetError(std::string("Connect failed: ") + StatusCodeString(status));
      return;
    }

    ctx->running.store(true);
    ctx->iterateThread = std::thread(ClientIterateLoop, ctx.get());

    {
      std::lock_guard<std::mutex> lock(g_clients_mutex);
      client_id_ = g_next_client_id.fetch_add(1);
      ctx->id = client_id_;
      g_clients[client_id_] = ctx;
    }
  }

  void OnOK() override {
    Napi::HandleScope scope(Env());
    deferred_.Resolve(Napi::Number::New(Env(), static_cast<double>(client_id_)));
  }

  void OnError(const Napi::Error& e) override {
    Napi::HandleScope scope(Env());
    deferred_.Reject(e.Value());
  }

 private:
  Napi::Promise::Deferred deferred_;
  ConnectionOptions options_;
  uint64_t client_id_;
  UA_StatusCode status_;
};

class DisconnectWorker : public Napi::AsyncWorker {
 public:
  DisconnectWorker(Napi::Promise::Deferred deferred, uint64_t client_id)
      : Napi::AsyncWorker(deferred.Env()), deferred_(deferred), client_id_(client_id) {}

  void Execute() override {
    DeleteClientContext(client_id_);
  }

  void OnOK() override {
    Napi::HandleScope scope(Env());
    deferred_.Resolve(Env().Undefined());
  }

  void OnError(const Napi::Error& e) override {
    Napi::HandleScope scope(Env());
    deferred_.Reject(e.Value());
  }

 private:
  Napi::Promise::Deferred deferred_;
  uint64_t client_id_;
};

class ReadWorker : public Napi::AsyncWorker {
 public:
  ReadWorker(Napi::Promise::Deferred deferred, std::shared_ptr<ClientContext> ctx,
             const std::string& node_id)
      : Napi::AsyncWorker(deferred.Env()),
        deferred_(deferred),
        ctx_(ctx),
        node_id_(node_id),
        status_(UA_STATUSCODE_GOOD) {
    UA_Variant_init(&value_);
  }

  ~ReadWorker() {
    UA_Variant_clear(&value_);
  }

  void Execute() override {
    UA_NodeId nodeId;
    if (!ParseNodeId(node_id_, &nodeId)) {
      SetError("Invalid NodeId: " + node_id_);
      return;
    }

    // Pause the iterate thread to avoid conflicting with synchronous API calls.
    PauseClientIteration(ctx_);

    {
      std::lock_guard<std::mutex> lock(ctx_->clientMutex);
      if (!ctx_->client) {
        UA_NodeId_clear(&nodeId);
        SetError("Client disconnected");
        ResumeClientIteration(ctx_);
        return;
      }

      status_ = UA_Client_readValueAttribute(ctx_->client, nodeId, &value_);
    }
    UA_NodeId_clear(&nodeId);

    ResumeClientIteration(ctx_);

    if (!UA_StatusCode_isGood(status_)) {
      SetError(std::string("Read failed: ") + StatusCodeString(status_));
      return;
    }
  }

  void OnOK() override {
    Napi::HandleScope scope(Env());
    Napi::Object result = Napi::Object::New(Env());
    result.Set("value", VariantToJsValue(Env(), &value_));
    result.Set("dataType", DataTypeName(value_.type));
    result.Set("statusCode", StatusCodeString(status_));
    deferred_.Resolve(result);
  }

  void OnError(const Napi::Error& e) override {
    Napi::HandleScope scope(Env());
    deferred_.Reject(e.Value());
  }

 private:
  Napi::Promise::Deferred deferred_;
  std::shared_ptr<ClientContext> ctx_;
  std::string node_id_;
  UA_Variant value_;
  UA_StatusCode status_;
};

class WriteWorker : public Napi::AsyncWorker {
 public:
  WriteWorker(Napi::Promise::Deferred deferred, std::shared_ptr<ClientContext> ctx,
              const std::string& node_id, const std::string& data_type, const Napi::Value& value)
      : Napi::AsyncWorker(deferred.Env()),
        deferred_(deferred),
        ctx_(ctx),
        node_id_(node_id),
        data_type_(data_type),
        status_(UA_STATUSCODE_GOOD) {
    UA_Variant_init(&valueVariant_);
    if (!JsValueToVariantScalar(Env(), value, data_type_, &valueVariant_)) {
      Napi::Error::New(Env(), "Failed to convert value to OPC UA type " + data_type_)
          .ThrowAsJavaScriptException();
    }
  }

  ~WriteWorker() {
    UA_Variant_clear(&valueVariant_);
  }

  void Execute() override {
    UA_NodeId nodeId;
    if (!ParseNodeId(node_id_, &nodeId)) {
      SetError("Invalid NodeId: " + node_id_);
      return;
    }

    // Pause the iterate thread to avoid conflicting with synchronous API calls.
    PauseClientIteration(ctx_);

    {
      std::lock_guard<std::mutex> lock(ctx_->clientMutex);
      if (!ctx_->client) {
        UA_NodeId_clear(&nodeId);
        SetError("Client disconnected");
        ResumeClientIteration(ctx_);
        return;
      }

      status_ = UA_Client_writeValueAttribute(ctx_->client, nodeId, &valueVariant_);
    }
    UA_NodeId_clear(&nodeId);

    ResumeClientIteration(ctx_);

    if (!UA_StatusCode_isGood(status_)) {
      SetError(std::string("Write failed: ") + StatusCodeString(status_));
      return;
    }
  }

  void OnOK() override {
    Napi::HandleScope scope(Env());
    Napi::Object result = Napi::Object::New(Env());
    result.Set("statusCode", StatusCodeString(status_));
    deferred_.Resolve(result);
  }

  void OnError(const Napi::Error& e) override {
    Napi::HandleScope scope(Env());
    deferred_.Reject(e.Value());
  }

 private:
  Napi::Promise::Deferred deferred_;
  std::shared_ptr<ClientContext> ctx_;
  std::string node_id_;
  std::string data_type_;
  UA_Variant valueVariant_;
  UA_StatusCode status_;
};

class SubscribeWorker : public Napi::AsyncWorker {
 public:
  SubscribeWorker(Napi::Promise::Deferred deferred, std::shared_ptr<ClientContext> ctx,
                  const std::string& node_id, double sampling_interval,
                  const Napi::Function& jsCallback)
      : Napi::AsyncWorker(deferred.Env()),
        deferred_(deferred),
        ctx_(ctx),
        node_id_(node_id),
        sampling_interval_(sampling_interval),
        subscription_id_(0),
        monitored_item_id_(0),
        status_(UA_STATUSCODE_GOOD) {
    tsfn_ = Napi::ThreadSafeFunction::New(
        Env(), jsCallback, "OpcUaSubscriptionCallback", 0, 1,
        [](Napi::Env) { /* finalize context optional */ });
  }

  void Execute() override {
    UA_NodeId nodeId;
    if (!ParseNodeId(node_id_, &nodeId)) {
      SetError("Invalid NodeId: " + node_id_);
      return;
    }

    // Pause the iterate thread to avoid conflicting with synchronous API calls.
    PauseClientIteration(ctx_);

    {
      std::lock_guard<std::mutex> lock(ctx_->clientMutex);
      if (!ctx_->client) {
        UA_NodeId_clear(&nodeId);
        SetError("Client disconnected");
        ResumeClientIteration(ctx_);
        return;
      }

      UA_CreateSubscriptionRequest subRequest = UA_CreateSubscriptionRequest_default();
      subRequest.requestedPublishingInterval = 500.0;
      UA_CreateSubscriptionResponse subResponse =
          UA_Client_Subscriptions_create(ctx_->client, subRequest, nullptr, nullptr, nullptr);
      if (!UA_StatusCode_isGood(subResponse.responseHeader.serviceResult)) {
        status_ = subResponse.responseHeader.serviceResult;
        UA_CreateSubscriptionResponse_clear(&subResponse);
        UA_NodeId_clear(&nodeId);
        SetError(std::string("Create subscription failed: ") + StatusCodeString(status_));
        ResumeClientIteration(ctx_);
        return;
      }
      subscription_id_ = subResponse.subscriptionId;
      UA_CreateSubscriptionResponse_clear(&subResponse);

      auto subCtx = std::make_unique<SubscriptionContext>();
      subCtx->tsfn = tsfn_;
      subCtx->nodeId = node_id_;
      subCtx->subscriptionId = subscription_id_;

      UA_MonitoredItemCreateRequest monRequest = UA_MonitoredItemCreateRequest_default(nodeId);
      monRequest.requestedParameters.samplingInterval = sampling_interval_;

      UA_MonitoredItemCreateResult monResponse = UA_Client_MonitoredItems_createDataChange(
          ctx_->client, subscription_id_, UA_TIMESTAMPSTORETURN_BOTH, monRequest,
          subCtx.get(), DataChangeCallback, nullptr);

      UA_NodeId_clear(&nodeId);

      if (!UA_StatusCode_isGood(monResponse.statusCode)) {
        status_ = monResponse.statusCode;
        UA_Client_Subscriptions_deleteSingle(ctx_->client, subscription_id_);
        SetError(std::string("Create monitored item failed: ") + StatusCodeString(status_));
        ResumeClientIteration(ctx_);
        return;
      }

      monitored_item_id_ = monResponse.monitoredItemId;
      subCtx->monitoredItemId = monitored_item_id_;
      ctx_->subscriptions[subscription_id_] = std::move(subCtx);
    }

    ResumeClientIteration(ctx_);
  }

  void OnOK() override {
    Napi::HandleScope scope(Env());
    Napi::Object result = Napi::Object::New(Env());
    result.Set("subscriptionId", subscription_id_);
    result.Set("monitoredItemId", monitored_item_id_);
    result.Set("statusCode", StatusCodeString(status_));
    deferred_.Resolve(result);
  }

  void OnError(const Napi::Error& e) override {
    tsfn_.Release();
    Napi::HandleScope scope(Env());
    deferred_.Reject(e.Value());
  }

 private:
  Napi::Promise::Deferred deferred_;
  std::shared_ptr<ClientContext> ctx_;
  std::string node_id_;
  double sampling_interval_;
  uint32_t subscription_id_;
  uint32_t monitored_item_id_;
  UA_StatusCode status_;
  Napi::ThreadSafeFunction tsfn_;
};

class UnsubscribeWorker : public Napi::AsyncWorker {
 public:
  UnsubscribeWorker(Napi::Promise::Deferred deferred, std::shared_ptr<ClientContext> ctx,
                    uint32_t subscription_id)
      : Napi::AsyncWorker(deferred.Env()),
        deferred_(deferred),
        ctx_(ctx),
        subscription_id_(subscription_id),
        status_(UA_STATUSCODE_GOOD) {}

  void Execute() override {
    // Pause the iterate thread to avoid conflicting with synchronous API calls.
    PauseClientIteration(ctx_);

    {
      std::lock_guard<std::mutex> lock(ctx_->clientMutex);
      if (ctx_->client) {
        status_ = UA_Client_Subscriptions_deleteSingle(ctx_->client, subscription_id_);
      }
    }

    auto it = ctx_->subscriptions.find(subscription_id_);
    if (it != ctx_->subscriptions.end()) {
      if (it->second && it->second->tsfn) {
        it->second->tsfn.Release();
      }
      ctx_->subscriptions.erase(it);
    }

    ResumeClientIteration(ctx_);

    if (!UA_StatusCode_isGood(status_)) {
      SetError(std::string("Unsubscribe failed: ") + StatusCodeString(status_));
    }
  }

  void OnOK() override {
    Napi::HandleScope scope(Env());
    Napi::Object result = Napi::Object::New(Env());
    result.Set("statusCode", StatusCodeString(status_));
    deferred_.Resolve(result);
  }

  void OnError(const Napi::Error& e) override {
    Napi::HandleScope scope(Env());
    deferred_.Reject(e.Value());
  }

 private:
  Napi::Promise::Deferred deferred_;
  std::shared_ptr<ClientContext> ctx_;
  uint32_t subscription_id_;
  UA_StatusCode status_;
};

// ---------------------------------------------------------------------------
// N-API method wrappers
// ---------------------------------------------------------------------------

Napi::Value ConnectAsync(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() < 1 || !info[0].IsObject()) {
    Napi::TypeError::New(env, "Expected options object").ThrowAsJavaScriptException();
    return env.Null();
  }

  ConnectionOptions options = ParseConnectionOptions(info[0].As<Napi::Object>());
  Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
  ConnectWorker* worker = new ConnectWorker(deferred, options);
  worker->Queue();
  return deferred.Promise();
}

Napi::Value DisconnectAsync(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() < 1 || !info[0].IsNumber()) {
    Napi::TypeError::New(env, "Expected clientId").ThrowAsJavaScriptException();
    return env.Null();
  }

  uint64_t client_id = static_cast<uint64_t>(info[0].As<Napi::Number>().Int64Value());
  Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
  DisconnectWorker* worker = new DisconnectWorker(deferred, client_id);
  worker->Queue();
  return deferred.Promise();
}

static std::shared_ptr<ClientContext> FindClientContext(uint64_t client_id) {
  std::lock_guard<std::mutex> lock(g_clients_mutex);
  auto it = g_clients.find(client_id);
  if (it == g_clients.end()) return nullptr;
  return it->second;
}

Napi::Value ReadAsync(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() < 2 || !info[0].IsNumber() || !info[1].IsString()) {
    Napi::TypeError::New(env, "Expected (clientId, nodeId)").ThrowAsJavaScriptException();
    return env.Null();
  }

  uint64_t client_id = static_cast<uint64_t>(info[0].As<Napi::Number>().Int64Value());
  std::string node_id = info[1].As<Napi::String>().Utf8Value();
  std::shared_ptr<ClientContext> ctx = FindClientContext(client_id);
  if (!ctx) {
    Napi::Error::New(env, "Client not found").ThrowAsJavaScriptException();
    return env.Null();
  }

  Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
  ReadWorker* worker = new ReadWorker(deferred, ctx, node_id);
  worker->Queue();
  return deferred.Promise();
}

Napi::Value WriteAsync(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() < 4 || !info[0].IsNumber() || !info[1].IsString() ||
      !info[3].IsString()) {
    Napi::TypeError::New(env, "Expected (clientId, nodeId, value, dataType)").ThrowAsJavaScriptException();
    return env.Null();
  }

  uint64_t client_id = static_cast<uint64_t>(info[0].As<Napi::Number>().Int64Value());
  std::string node_id = info[1].As<Napi::String>().Utf8Value();
  std::string data_type = info[3].As<Napi::String>().Utf8Value();
  std::shared_ptr<ClientContext> ctx = FindClientContext(client_id);
  if (!ctx) {
    Napi::Error::New(env, "Client not found").ThrowAsJavaScriptException();
    return env.Null();
  }

  Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
  WriteWorker* worker = new WriteWorker(deferred, ctx, node_id, data_type, info[2]);
  worker->Queue();
  return deferred.Promise();
}

Napi::Value SubscribeAsync(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() < 4 || !info[0].IsNumber() || !info[1].IsString() ||
      !info[2].IsNumber() || !info[3].IsFunction()) {
    Napi::TypeError::New(env, "Expected (clientId, nodeId, samplingInterval, callback)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }

  uint64_t client_id = static_cast<uint64_t>(info[0].As<Napi::Number>().Int64Value());
  std::string node_id = info[1].As<Napi::String>().Utf8Value();
  double sampling_interval = info[2].As<Napi::Number>().DoubleValue();
  Napi::Function callback = info[3].As<Napi::Function>();
  std::shared_ptr<ClientContext> ctx = FindClientContext(client_id);
  if (!ctx) {
    Napi::Error::New(env, "Client not found").ThrowAsJavaScriptException();
    return env.Null();
  }

  Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
  SubscribeWorker* worker =
      new SubscribeWorker(deferred, ctx, node_id, sampling_interval, callback);
  worker->Queue();
  return deferred.Promise();
}

Napi::Value UnsubscribeAsync(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() < 2 || !info[0].IsNumber() || !info[1].IsNumber()) {
    Napi::TypeError::New(env, "Expected (clientId, subscriptionId)").ThrowAsJavaScriptException();
    return env.Null();
  }

  uint64_t client_id = static_cast<uint64_t>(info[0].As<Napi::Number>().Int64Value());
  uint32_t subscription_id = info[1].As<Napi::Number>().Uint32Value();
  std::shared_ptr<ClientContext> ctx = FindClientContext(client_id);
  if (!ctx) {
    Napi::Error::New(env, "Client not found").ThrowAsJavaScriptException();
    return env.Null();
  }

  Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
  UnsubscribeWorker* worker = new UnsubscribeWorker(deferred, ctx, subscription_id);
  worker->Queue();
  return deferred.Promise();
}

}  // namespace

Napi::Object Init(Napi::Env env, Napi::Object exports) {
  exports.Set("connectAsync", Napi::Function::New(env, ConnectAsync));
  exports.Set("disconnectAsync", Napi::Function::New(env, DisconnectAsync));
  exports.Set("readAsync", Napi::Function::New(env, ReadAsync));
  exports.Set("writeAsync", Napi::Function::New(env, WriteAsync));
  exports.Set("subscribeAsync", Napi::Function::New(env, SubscribeAsync));
  exports.Set("unsubscribeAsync", Napi::Function::New(env, UnsubscribeAsync));
  return exports;
}

NODE_API_MODULE(opcua_open62541, Init)
