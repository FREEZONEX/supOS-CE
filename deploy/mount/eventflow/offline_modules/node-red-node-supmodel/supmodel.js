const axios = require('axios');
const fs = require('fs');
const path = require('path');

module.exports = function (RED) {

    function SelectModel(config) {
        // 模型定义数据
        let restQueryPageStartKey = '';
        let restQueryPageStartValue = '0';

        RED.nodes.createNode(this, config);

        const node = this;

        node.intervalId="";
        node.selectedModel = config.selectedModel || "";
        node.modelSchemaObj = null;

        if (config.modelSchema) {
            node.modelSchemaObj = JSON.parse(config.modelSchema);
        }
        node.modelMapping = null;
        if (config.modelMapping) {
            node.modelMapping = JSON.parse(config.modelMapping);
        }

        if (node.selectedModel) {
            // console.log("1================================================>>>\n")
            // console.log(node.modelSchema);
            // console.log("<<<================================================1\n")

            var protocol = "";
            if (node.modelSchemaObj?.protocol?.protocol) {
                protocol = node.modelSchemaObj.protocol.protocol;
            }
            console.log("set protocol value: " + protocol + ", topic: " + node.selectedModel);
      
            // 模型数据缓存不存在需要从服务端重新拉一份
            if (!node.modelSchemaObj && node.selectedModel != 'Auto') {
                queryModelSchema(node.selectedModel).then((res) => {
                    if (res?.data?.data) {
                        node.modelSchemaObj = res.data.data;
                    } else {
                        console.log(`topic ${node.selectedModel} query failed`);
                    }
                });
            }

            if (node.intervalId) {
                clearInterval(node.intervalId);
            }
            node.intervalId = triggerDeleteLogFile("/data/log/topology.log");

            node.on('input', function (msg) {
                var modelTopic = msg.model || (node.selectedModel == 'Auto' ? '' : node.selectedModel);
                
                var transferred = protocolDataTransfer(node.modelSchemaObj, msg, modelTopic);
                if (transferred.code == 200) {

                    if (transferred.isMultiple === true) {
                        node.send([transferred.data, null]);
                    } else {
                        msg.topic = rebuildSendTopic(modelTopic);
                        msg.payload = transferred.data;
                        node.send([msg, null]);
                    }
                    
                } else {

                    if (transferred.code >= 200 && transferred.code < 300) {
                        node.debug(transferred.message);
                    } else if (transferred.code >= 300 && transferred.code < 400) { // 错误写到日志文件，拓扑图分析
                        let topics = [];
                        if (node.selectedModel && node.selectedModel != 'Auto') {
                            topics.push(node.selectedModel);
                        } else if (node.modelMapping) {
                            topics = [...new Set(Object.values(node.modelMapping))];
                        }
                        if (topics.length > 0) {
                            try {
                                var filePath = "/data/log/topology.log";
                                if (transferred.code == 301) {
                                    writeToFile(topics, '', filePath); // 恢复
                                } else {
                                    writeToFile(topics, transferred.message, filePath);
                                }
                            } catch (e) {
                                console.error('write error log fail:', e.message);
                            }
                        }
                    } else {
                        node.error(transferred.message, msg);
                    }
                }
            });
            
        }

        function triggerDeleteLogFile(filePath) {
            var intervalId = setInterval(function() {
                fs.stat(filePath, (err, stats) => {
                    if (err) {
                        return;
                    }
                
                    // 获取文件的创建时间（或修改时间）
                    const fileCreationTime = stats.birthtime || stats.mtime; // birthtime 是创建时间，mtime 是修改时间
                    
                    const currentTime = new Date().getTime();

                    // 计算时间差（单位：毫秒）
                    const timeDifference = currentTime - fileCreationTime;
                    const oneHourInMilliseconds = 60 * 60 * 1000; // 1 小时的毫秒数
                
                    // 判断是否超过 1 小时
                    if (timeDifference > oneHourInMilliseconds) {
                        // 删除文件
                        fs.unlink(filePath, (err) => {
                            if (err) {
                                console.error('删除文件失败:', err);
                            } else {
                                console.log('文件已删除:', filePath);
                            }
                        });
                    } 
                });
            }, 60 * 60 * 1000); // 1小时删除一次文件
            return intervalId;
        }

        function writeToFile(topics, errorMsg, filePath) {
            let errorJsons = "";
            for (let i in topics) {
                var tts = topics[i];
                if (Array.isArray(tts)) {
                    for (let j = 0; j < tts.length; j++) {
                        var topic = tts[j].split(":")[0];
                        var current = new Date().getTime();
                        if (errorMsg) {
                            let errJson = `{"instanceTopic":"${topic}","topologyNode":"pushOriginalData","eventCode":"1","eventMessage":"${errorMsg}","eventTime":${current}},\n`
                            errorJsons += errJson;
                        } else {
                            let errJson = `{"instanceTopic":"${topic}","topologyNode":"pushOriginalData","eventCode":"0","eventMessage":"","eventTime":${current}},\n`
                            errorJsons += errJson;
                        }
                    }
                } else {
                    var topic = tts.split(":")[0];
                    var current = new Date().getTime();
                    if (errorMsg) {
                        let errJson = `{"instanceTopic":"${topic}","topologyNode":"pushOriginalData","eventCode":"1","eventMessage":"${errorMsg}","eventTime":${current}},\n`
                        errorJsons += errJson;
                    } else {
                        let errJson = `{"instanceTopic":"${topic}","topologyNode":"pushOriginalData","eventCode":"0","eventMessage":"","eventTime":${current}},\n`
                        errorJsons += errJson;
                    }
                }
            }
            fs.appendFile(filePath, errorJsons, (error) => {
                if (error) {
                    console.log(error);
                }
            });
        }

        function queryModelSchema(topic) {
            var encodeUrl = encodeURI('http://uns:8080/inter-api/supos/uns/instance?topic=' + topic);
            return axios.get(encodeUrl).then(res => {
                return res;
            }).catch(e => console.log(e.message));
        }

        function protocolDataTransfer(modelDef, msg, modelTopic) {
            var protocol = modelDef?.protocol?.protocol || "relation";
            switch(protocol) {
                case "modbus": 
                    return arrayDataTransfer(modelTopic, msg.payload, node.modelMapping);
                case "opcda": 
                    return opcdaDataTransfer(msg.payload, node.modelMapping);
                case "opcua": {
                    if (msg.topic == 'subscriptionId') {
                        return {
                            "code": 205,
                            "message": 'subscription event triggered'
                        }
                    } 
                    if (msg.status) {
                        let liveStatus = ['keepalive', 'connected', 'subscribed'];
                        node.debug("receive opcua server status is " + msg.status);
                        if (msg.status.indexOf('active') !== -1 || msg.status.indexOf('connected') !== -1 || liveStatus.includes(msg.status)) {
                            return {
                                "code": 301,
                                "message": 'opcua server is connected.'
                            }
                        } else if (msg.status == 'create client' || msg.status == 'connecting' ) {
                            return {
                                "code": 205,
                                "message": 'opcua client is connecting...'
                            }
                        } else {
                            return {
                                "code": 302,
                                "message": 'opcua server connect failed, status: ' + msg.status
                            }
                        }
                    }
                    if (msg.payload != null && msg.payload != undefined) {
                        return opcuaDataTransfer(node.modelMapping, msg);
                    } else {
                        return {
                            "code": 400,
                            "message": `opcua message payload(${msg.payload}) is empty.`
                        }
                    }
                }
                case "rest": 
                    if (msg.payload != null && msg.payload != undefined) {
                        return restDataTransfer(msg.payload, modelDef);
                    } else {
                        return {
                            "code": 400,
                            "message": `rest message payload(${msg.payload}) is empty.`
                        }
                }
                case "icmp": 
                    return icmpDataTransfer(msg.payload);
                case "relation": case "mqtt": 
                    return mqttDataTransfer(msg.payload);
                default: {
                    // 目前支持3种数据类型转换，分别是数组、基本数据类型、json对象，前面2种可以分别按照modbus和opcua处理
                    var inputPayload = msg.payload;
                    if (Array.isArray(inputPayload)) {
                        return arrayDataTransfer(modelTopic, inputPayload, node.modelMapping);
                    } else if (typeof inputPayload === 'object') {
                        return jsonDataTransfer(node.modelMapping, inputPayload);
                    } else if (isPrimitive(inputPayload)) {
                        return arrayDataTransfer(modelTopic, [inputPayload], node.modelMapping);
                    } else {
                        return {
                            "code": 400,
                            "message": `自定义协议目前只支持数组、JSON数据类型.`
                        }
                    }
                }
            }
        };

        // 判断是否为基础数据类型
        function isPrimitive(value) {
            return (
              typeof value === 'string' ||
              typeof value === 'number' ||
              typeof value === 'boolean' ||
              typeof value === 'symbol' ||
              value === null
            );
        }

        function mqttDataTransfer(inputData) {
            return {
                "code": 200,
                "isMultiple": false,
                "data": {
                    "_resource_": inputData
                }
            }
        }

        function icmpDataTransfer(inputStatus) {
            let fixTime = new Date().getTime();
            var data = {
                "_resource_": {
                    "_ct": fixTime,
                    "_qos": 0,
                    "status": inputStatus === true ? 1 : 0
                }
            }
            return {
                "code": 200,
                "isMultiple": false,
                "data": data
            }
        }

        // 康熙说只传原始数据，不要包装
        function restDataTransfer(inputData, modelDef) {
            // 最外层数据路径， 例如： 输入{"data": {"list": []}}, 那么dataPath=data.list
            var dataPath = modelDef['dataPath'];
            var data = inputData;
            if (typeof inputData === 'string') {
                try {
                    data = JSON.parse(inputData);
                } catch(err) {
                    return {
                        "code": 400,
                        "message": inputData
                    }
                }
            }
            if (dataPath) {
                var indexs = dataPath.split('.');
                for (let i in indexs) {
                    data = data[indexs[i]];
                }
            }
            if (!data) {
                return {
                    "code": 400,
                    "message": `invalid data ${inputData}`
                }
            }
            // 数据参数映射
            var newData = [];
            var mappingCache = {};
            if (Array.isArray(data)) {
                // 对数组内每个对象的属性根据模型映射关系进行匹配
                for (let obj of data) {
                    var newObj = getMappingField(obj, mappingCache, modelDef);
                    newData.push(newObj);
                }
            } else {
                var newObj = getMappingField(data, mappingCache, modelDef);
                newData.push(newObj);
            }
            if (newData.length == 0) {
                node.context().flow.set("pageStart", 0);
            }
            // 更新页码, 以备下一次请求使用
            if (restQueryPageStartKey) {
                updateUrl(newData.length);
            }
            return {
                "code": 200,
                "isMultiple": false,
                "data": {
                    "_source_": inputData,
                    "_resource_": newData
                }
            }
            // return {
            //     "code": 200,
            //     "isMultiple": false,
            //     "data": inputData 
            // }
        }

        function isEmptyObject(obj) {
            return Object.keys(obj).length === 0;
        }
        
        function getMappingField(inputObj, mapping, modelDef) {
            var modelFields = modelDef['fields'];
            var newObj = {};
            var isEmpty = isEmptyObject(mapping);
            for (let [key, value] of Object.entries(inputObj)) {
                if (isEmpty) {
                    for (var field of modelFields) {
                        if (field['index'] == key) {
                            newObj[field['name']] = value;
                            mapping[key] = field['name'];
                            break;
                        }
                    }
                } else {
                    var fieldName = mapping[key];
                    if (fieldName) {
                        newObj[fieldName] = value;
                    }
                }
            }
            return newObj;
        }

        function jsonDataTransfer(mappings, inputPayload) {
            var multipleVqt = [];
            var timestamp = Date.now();
            let resultMap = {};  // {topic1: [{field1: value1}, {field2: value2}]}

            // {位号1: value1, 位号2: value2} ==> {topic1: [{field1: value1}, {field2: value2}]}
            for (let key in inputPayload) {

                var ref = mappings[key]; // mappings结构为 {位号: [unsTopic:fieldName]}

                if (!Array.isArray(ref)) {
                    ref = [ref]; // [unsTopic:fieldName]
                } 

                for (let i in ref) {
                    let splitArr = ref[i].split(":");
                    let fieldName = splitArr[1] || "v";
                    let modelTopic = splitArr[0];

                    var fieldValues = resultMap[modelTopic];
                    if (!fieldValues) {
                        fieldValues = [];
                    }
                    var single = {};
                    single[fieldName] = inputPayload[key];
                    fieldValues.push(single);

                    resultMap[modelTopic] = fieldValues
                }
                
            }

            for (let tp in resultMap) {
                var valueArray = resultMap[tp]; // [{field1: value1}, {field2: value2}]
                let vqt = {
                    "_qos": 0,
                    "_ct": timestamp,
                }
                for (let v in valueArray) {
                    for (let f in valueArray[v]) { // valueArray[v] => {field1: value1}
                        vqt[f] = valueArray[v][f];
                    }
                }

                multipleVqt.push({
                    "topic": rebuildSendTopic(tp),
                    "payload": {
                        "_source_": inputPayload,
                        "_resource_": vqt
                    }
                });
            }
            return {
                "code": 200,
                "isMultiple": true,
                "data": multipleVqt
            };
        }

        // data is array of number
        function arrayDataTransfer(modelTopic, data, mappings) {
            if (!Array.isArray(data))  {
                var errMsg = RED._("supmodel.protocol_data_not_match");
                return {
                    "code": 400,
                    "message": errMsg
                };
            }
            var unsData;
            var timestamp = Date.now();

            if (mappings) {
                unsData = {
                    "_qos": 0,
                    "_ct": timestamp
                };
                let indexMap = mappings[modelTopic];
                if (indexMap) {
                    for (let k in indexMap) { // key是属性名，value是数组下标
                        var realVal = data[indexMap[k]];
                        if (realVal != null && realVal != undefined) {
                            unsData[k] = realVal;
                        }
                    }
                }
            } else {
                node.error(`${modelTopic} mapping cache is null`);
                return {
                    "code": 400,
                    "message": `${modelTopic} mapping cache is null`
                };
            }
            return {
                "code": 200,
                "isMultiple": false,
                "data": {
                    "_source_": data,
                    "_resource_": unsData
                }
            };
        };

        function opcdaDataTransfer(inputArrayData, mappings) {
            let multipleVqt = [];
            let vqtMap = {};

            for (let i in inputArrayData) {
                let nodeId = inputArrayData[i].itemID; // opcda的位号
                let fields = mappings[nodeId];
                if (fields) {
                    for (let j in fields) {
                        var vqt = {
                            "_qos": inputArrayData[i].errorCode,
                            "_ct": new Date(inputArrayData[i].timestamp).getTime()
                        }
                        let field = fields[j]; // topic:fieldName
                        let sp = field.split(':');
                        var fieldName = sp[1];
                        vqt[fieldName] = inputArrayData[i].value;
                        var topic = sp[0];
                        if (!vqtMap[topic]) {
                            vqtMap[topic] = [];
                        }
                        vqtMap[topic].push(vqt);
                        
                    }

                } else {
                    console.log(`opcda ${nodeId} 映射uns字段不存在`);
                }
            }
            // 根据topic进行合并
            for (let k in vqtMap) {
                multipleVqt.push({
                    "topic": rebuildSendTopic(k),
                    "payload": {
                        "_source_": inputArrayData,
                        "_resource_": vqtMap[k]
                    }
                });
            }
            
            return {
                "code": 200,
                "isMultiple": true,
                "data": multipleVqt
            };
        }

        function rebuildSendTopic(originTopic) {
            return originTopic;
        }

        function opcuaDataTransfer(mappings, msg) {
            var vqtWrapper = {};
            // 判断是单个订阅还是批量订阅
            let batch = (msg.payload?.value?.value != null && msg.payload?.value?.value != undefined) ? true : false;
            // 寻找opcua位号对应uns的topic, msg.topic为opcua位号地址
            let fts = [];
            if (mappings && mappings[msg.topic]) {
                let ts = mappings[msg.topic]; // 0=uns topic, 1=fieldName
                if (!Array.isArray(ts)) {
                    let tsArr = [];
                    tsArr.push(ts); // 兼容老数据
                    ts = tsArr;
                }
                for (let i in ts) {
                    let sp = ts[i].split(':');
                    fts.push({
                        "modelTopic": sp[0],
                        "fieldName": sp[1]
                    });
                }
            }
            
            let multipleVqt = [];
            if (batch) {
                if (fts.length > 0) {
                    for (let j in fts) {
                        var vqt = {
                            "_qos": msg.payload.statusCode ? msg.payload.statusCode._value : -1,
                            "_ct": new Date(msg.payload.sourceTimestamp).getTime()
                        }

                        var fname = fts[j].fieldName;
                        vqt[fname] = msg.payload.value.value;

                        multipleVqt.push({
                            "topic": rebuildSendTopic(fts[j].modelTopic),
                            "payload": {
                                "_source_": msg.payload,
                                "_resource_": vqt
                            }
                        });
                    }
                } else {
                    return {
                        "code": 400,
                        "message": `fields mapping cache is null`
                    };
                }
                
            } else {
                if (fts.length > 0) {
                    for (let j in fts) {
                        var vqt = {
                            "_qos": msg.statusCode ? msg.statusCode._value : -1,
                            "_ct": new Date(msg.sourceTimestamp).getTime()
                        }
                        var fname = fts[j].fieldName;
                        vqt[fname] = msg.payload;

                        multipleVqt.push({
                            "topic": rebuildSendTopic(fts[j].modelTopic),
                            "payload": {
                                "_source_": msg.payload,
                                "_resource_": vqt
                            }
                        });
                    }
                } else {
                    return {
                        "code": 400,
                        "message": `fields mapping cache is null`
                    };
                }
            }
            
            return {
                "code": 200,
                "isMultiple": true,
                "data": multipleVqt
            };
    
        };

        function singleDataTransfer(mappings, msg) {
            var fieldName = "v";
            msg.model = msg.topic;

            if (mappings && mappings[msg.topic]) {
                let ts = mappings[msg.topic].split(':'); // 0=uns topic, 1=fieldName
                msg.model = ts[0];
                fieldName = ts[1] || "v";
            }
            let vqt = {
                "_qos": 0,
                "_ct": new Date().getTime()
            }
            vqt[fieldName] = msg.payload;

            return {
                "code": 200,
                "isMultiple": false,
                "data": {
                    "_source_": msg.payload,
                    "_resource_": vqt
                }
            };
    
        };
    }
    RED.nodes.registerType("supmodel", SelectModel);
    
    
}