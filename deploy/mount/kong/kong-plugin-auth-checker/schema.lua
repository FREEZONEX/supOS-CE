local typedefs = require "kong.db.schema.typedefs"

return {
  name = "auth-checker",
  fields = {
    -- 配置项: auth_url, 指定 /auth 请求的 URL
    { config = {
        type = "record",
        fields = {
          { whitelist_paths = {  -- 需要放行的路径列表
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
                  "^/files.*$",
                  "^/tier0-login.*$",
                  "^/nodered.*$",
                  "^/inter-api/supos/cascade.*$",
                  "^/inter-api/supos/license.*$",
                  "^/swagger-ui.*$",
                  "/403"
              },
          },
          },
          -- 是否启用资源权限检查
          { enable_resource_check = { type = "boolean", default = false } },
          -- 是否启用拒绝策略资源权限检查
          { enable_deny_check = { type = "boolean", default = true } },
          -- 登录页URL
          { login_url = { type = "string", default = "/tier0-login"} },
          -- 403页面URL
          { forbidden_url = { type = "string", default = "/403"} },
        },
      },
    },
  },
}