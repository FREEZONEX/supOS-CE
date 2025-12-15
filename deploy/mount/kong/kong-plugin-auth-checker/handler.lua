local http = require "resty.http"
local ngx = ngx
local cjson = require "cjson"

local plugin = {
  PRIORITY = 1000,  -- 插件的优先级
  VERSION = "1.0.0",  -- 插件版本
}

local auth_url = "http://uns:8080/inter-api/supos/auth/userinfo"

function isNotEmpty(str)
  return str and str ~= "" and str ~= "null"
end

-- 判断路径是否在白名单中
local function is_whitelisted_path(path, whitelist_paths)
  for _, whitelist_path in ipairs(whitelist_paths or {}) do
    path = path:gsub("-", "")
    if string.match(path, whitelist_path) then
      return true
    end
  end
  return false
end

local function unauthorized(redirect_url)
  --ngx.log(ngx.ERR, "重定向地址》》》: ", redirect_url)
  if not redirect_url or redirect_url == "" then
    return ngx.exit(ngx.HTTP_UNAUTHORIZED)
  else
    --ngx.log(ngx.ERR, "重定向到: ", redirect_url)
    -- 设置响应头以清空浏览器的 Cookie
    ngx.header["Set-Cookie"] = {
      "supos_community_token=; Path=/; Max-Age=0;"
    }
    return ngx.redirect(redirect_url, ngx.HTTP_MOVED_TEMPORARILY)
  end
end

local function forbidden(redirect_url)
  if not redirect_url or redirect_url == "" then
    return ngx.exit(ngx.HTTP_FORBIDDEN)
  else
    --ngx.log(ngx.ERR, "重定向到: ", redirect_url)
    return ngx.redirect(redirect_url, ngx.HTTP_MOVED_TEMPORARILY)
  end
end

-- 定义函数来判断 currentMethod 是否在 methods 数组中 true存在，false不存在
local function isMethodMatched(methodList, currentMethod)
  -- 检查methodList是否为空或是空数组
  if not methodList or #methodList == 0 then
    return true
  end

  -- 遍历methodList，检查是否有匹配的method
  for _, method in ipairs(methodList) do
    if method == currentMethod then
      return true
    end
  end

  return false
end


-- 处理请求的访问逻辑
function plugin:access(conf)
  local current_path = ngx.var.uri  -- 获取当前请求路径
  local current_method = string.lower(ngx.req.get_method()) -- 获取当前请求方法
  --ngx.log(ngx.ERR, ">>>>>>>>>>>>当前请求方法：", current_method, " 当前请求路径：", current_path)

  -- 检查是否在白名单路径中
  if is_whitelisted_path(current_path, conf.whitelist_paths) then
    --ngx.log(ngx.ERR, ">>>>>>>>>>>>当前请求路径： ", current_path, " 是白名单内路径，放行")
    return  -- 直接放行
  end

  -- 获取所有 cookies
  local cookies = ngx.req.get_headers()["Cookie"]

  if not cookies then
    --ngx.log(ngx.ERR, ">>>>>>>>>>>>没有找到Cookie，拒绝本次请求")
    --return ngx.exit(ngx.HTTP_UNAUTHORIZED)  -- 如果没有 cookies，返回 401
    return unauthorized(conf.login_url)
  end

  local supos_community_token = string.match(cookies, "supos_community_token=([^;]+)") -- 提取 supos_community_token 的值
  --ngx.log(ngx.ERR, ">>>>>>>>>>>>supos_community_token:",supos_community_token)

  -- 如果未找到 supos_community_token，返回 401
  if not supos_community_token then
    --ngx.log(ngx.ERR, ">>>>>>>>>>>>没有找到supos_community_token，拒绝本次请求")
    --return ngx.exit(ngx.HTTP_UNAUTHORIZED)  -- 如果没有 cookies，返回 401
    return unauthorized(conf.login_url)
  end

  -- 向后端发送 POST 请求，带上所有 cookies
  local httpc = http.new()
  local res, err = httpc:request_uri(auth_url, {
    method = "GET",
    headers = {
      ["Cookie"] = "supos_community_token=" .. supos_community_token, -- 携带 supos_community_token
    },
  })

  if not res or res.status ~= 200 then
    ngx.log(ngx.ERR, ">>>>>>> 请求backend接口失败: ", err)
    --return ngx.exit(ngx.HTTP_UNAUTHORIZED)  -- 如果请求失败，返回 401
    return unauthorized(conf.login_url)
  end

  local success, json_data = pcall(cjson.decode, res.body)

  if not success then
    ngx.log(ngx.ERR, "解析backend返回JSON失败 response: ", json_data)
    return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
  end

  if json_data.data.superAdmin == true then
    ngx.log(ngx.INFO, "当前用户为超级管理员，直接放行")
    return
  end

    ------------------拒绝策略--------------------------
  --ngx.log(ngx.ERR, ">>>>>>>>>>>>>>资源权限校验开关: ", conf.enable_resource_check)
  if conf.enable_deny_check then
    -- 解析接口返回的 JSON 数据
    local success, json_data = pcall(cjson.decode, res.body)

    if not success then
      ngx.log(ngx.ERR, "解析backend返回JSON失败 response: ", json_data)
      return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    -- 判断请求路径是否在 deny_paths 中
    local is_deny_uri = false
    local is_deny_method = false

    if json_data.data.denyResourceList and #json_data.data.denyResourceList > 0 then
      -- 提取 denyResourceList 中的 uri 并构建 deny_paths 数组
      for _, resource in ipairs(json_data.data.denyResourceList) do
        local lower_deny_path = string.lower(resource.uri)
        current_path = current_path:gsub("-", "")
        current_path = string.lower(current_path)
        -- 如果请求路径与某个允许的路径匹配，正则匹配
        if string.match(current_path, "^" .. lower_deny_path .. ".*$") then
          is_deny_uri = true
          is_deny_method = isMethodMatched(resource.methods,current_method)
          break
        end
      end
      --ngx.log(ngx.ERR, "用户拒绝访问的路径为： ", cjson.encode(deny_paths))


      -- 如果请求路径匹配，执行拒绝策略 返回403
      if is_deny_uri and is_deny_method then
        ngx.log(ngx.ERR, "路径 " .. current_path .. " 命中拒绝策略，当前请求被拒绝")
        --return ngx.exit(ngx.HTTP_FORBIDDEN)
        return forbidden(conf.forbidden_url)
      --else
      --  --放行
      --  --ngx.log(ngx.ERR, "路径 " .. current_path .. " 存在，允许放行")
      --  return
      end
    end
  end
  ------------------拒绝策略END--------------------------

  ------------------资源权限校验--------------------------
  --ngx.log(ngx.ERR, ">>>>>>>>>>>>>>资源权限校验开关: ", conf.enable_resource_check)
  if conf.enable_resource_check then
    -- 解析接口返回的 JSON 数据（用户有权限访问的路径数组）
    local success, json_data = pcall(cjson.decode, res.body)

    if not success then
      ngx.log(ngx.ERR, "解析backend返回JSON失败 response: ", json_data)
      return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    if json_data.data.resourceList and #json_data.data.resourceList > 0 then
      -- 判断请求路径是否在 allowed_paths 中
      local is_allowed_uri = false
      local is_allowed_method = false
      -- 提取 resourceList 中的 uri 并构建 allowed_paths 数组
      for _, resource in ipairs(json_data.data.resourceList) do
        local lower_allowed_path = string.lower(resource.uri)
        current_path = current_path:gsub("-", "")
        current_path = string.lower(current_path)
        --ngx.log(ngx.ERR, "当前允许路径： " .. lower_allowed_path .. "当前路径：" ..current_path)
        -- 如果请求路径与某个允许的路径匹配，正则匹配
        if string.match(current_path, "^" .. lower_allowed_path .. ".*$") then
          is_allowed_uri = true
          is_allowed_method = isMethodMatched(resource.methods,current_method)
          break
        end
      end
      --ngx.log(ngx.ERR, "用户允许访问的路径为：", cjson.encode(allowed_paths))

      -- 如果请求路径匹配，放行请求；否则返回 401
      if is_allowed_uri and is_allowed_method then
        --ngx.log(ngx.ERR, "路径 " .. current_path .. " 存在，允许放行")
        return
      else
        ngx.log(ngx.ERR, "路径 " .. current_path .. " 资源检查不存在，拒绝放行")
        --return ngx.exit(ngx.HTTP_FORBIDDEN)
        return forbidden(conf.forbidden_url)
      end
    end

  end
  --ngx.log(ngx.INFO, ">>>>>>>>>验证access token成功，继续执行原始请求")
  ------------------资源权限校验END--------------------------
end

return plugin