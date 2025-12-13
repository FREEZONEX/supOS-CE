local typedefs = require "kong.db.schema.typedefs"

return {
  name = "auth-checker",
  fields = {
    { config = {
        type = "record",
        fields = {
          { whitelist_paths = {
              type = "array",
              elements = { type = "string" },
              default = {
                  "^/inter-api/supos/auth.*$",
                  "^/inter-api/supos/systemConfig.*$",
                  "^/$",
                  "^/assets.*$",
                  "^/keycloak.*$",
                  "^/locale.*$",
                  "^/gitea.*git.*$",
                  "^/logo.*$",
                  "^/tier0-login.*$",
                  "/403"
              },
          },
          },
          { enable_resource_check = { type = "boolean", default = false } },
          { enable_deny_check = { type = "boolean", default = true } },
          { login_url = { type = "string", default = "/tier0-login"} },
          { forbidden_url = { type = "string", default = "/403"} },
        },
      },
    },
  },
}

