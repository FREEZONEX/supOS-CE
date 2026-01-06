local http = require "resty.http"
local ngx = ngx
local cjson = require "cjson"

local plugin = {
  PRIORITY = 1000,
  VERSION = "1.0.0",
}

local auth_url = "http://uns:8080/inter-api/supos/auth/userinfo"

function isNotEmpty(str)
  return str and str ~= "" and str ~= "null"
end

-- Check if path is in whitelist
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
  if not redirect_url or redirect_url == "" then
    return ngx.exit(ngx.HTTP_UNAUTHORIZED)
  else
    -- Clear browser cookie
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
    return ngx.redirect(redirect_url, ngx.HTTP_MOVED_TEMPORARILY)
  end
end

-- Check if currentMethod is in methods array
local function isMethodMatched(methodList, currentMethod)
  -- If methodList is empty or nil, match all methods
  if not methodList or #methodList == 0 then
    return true
  end

  -- Check if currentMethod is in methodList
  for _, method in ipairs(methodList) do
    if method == currentMethod then
      return true
    end
  end

  return false
end


-- Handle request access logic
function plugin:access(conf)
  local current_path = ngx.var.uri
  local current_method = string.lower(ngx.req.get_method())

  -- Check if path is in whitelist
  if is_whitelisted_path(current_path, conf.whitelist_paths) then
    return
  end

  -- Get all cookies
  local cookies = ngx.req.get_headers()["Cookie"]

  if not cookies then
    return unauthorized(conf.login_url)
  end

  local supos_community_token = string.match(cookies, "supos_community_token=([^;]+)")

  -- If supos_community_token not found, return 401
  if not supos_community_token then
    return unauthorized(conf.login_url)
  end

  -- Send GET request to backend with supos_community_token
  local httpc = http.new()
  local res, err = httpc:request_uri(auth_url, {
    method = "GET",
    headers = {
      ["Cookie"] = "supos_community_token=" .. supos_community_token,
    },
  })

  if not res or res.status ~= 200 then
    ngx.log(ngx.ERR, ">>>>>>> Request to backend failed: ", err)
    return unauthorized(conf.login_url)
  end

  -- Deny policy check
  if conf.enable_deny_check then
    -- Parse JSON response from backend
    local success, json_data = pcall(cjson.decode, res.body)

    if not success then
      ngx.log(ngx.ERR, "Failed to parse backend JSON response: ", json_data)
      return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    -- Check if request path is in deny_paths
    local is_deny_uri = false
    local is_deny_method = false

    if json_data.data.denyResourceList and #json_data.data.denyResourceList > 0 then
      -- Build deny_paths array from denyResourceList
      for _, resource in ipairs(json_data.data.denyResourceList) do
        local lower_deny_path = string.lower(resource.uri)
        current_path = current_path:gsub("-", "")
        current_path = string.lower(current_path)
        -- Match request path with deny path using regex
        if string.match(current_path, "^" .. lower_deny_path .. ".*$") then
          is_deny_uri = true
          is_deny_method = isMethodMatched(resource.methods, current_method)
          break
        end
      end

      -- If path matches deny policy, return 403
      if is_deny_uri and is_deny_method then
        ngx.log(ngx.ERR, "Path " .. current_path .. " matches deny policy, request denied")
        return forbidden(conf.forbidden_url)
      end
    end
  end

  -- Resource permission check
  if conf.enable_resource_check then
    -- Parse JSON response (paths user has permission to access)
    local success, json_data = pcall(cjson.decode, res.body)

    if not success then
      ngx.log(ngx.ERR, "Failed to parse backend JSON response: ", json_data)
      return ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    if json_data.data.resourceList and #json_data.data.resourceList > 0 then
      -- Check if request path is in allowed_paths
      local is_allowed_uri = false
      local is_allowed_method = false
      -- Build allowed_paths array from resourceList
      for _, resource in ipairs(json_data.data.resourceList) do
        local lower_allowed_path = string.lower(resource.uri)
        current_path = current_path:gsub("-", "")
        current_path = string.lower(current_path)
        -- Match request path with allowed path using regex
        if string.match(current_path, "^" .. lower_allowed_path .. ".*$") then
          is_allowed_uri = true
          is_allowed_method = isMethodMatched(resource.methods, current_method)
          break
        end
      end

      -- If path matches, allow request; otherwise return 403
      if is_allowed_uri and is_allowed_method then
        return
      else
        ngx.log(ngx.ERR, "Path " .. current_path .. " failed resource check, request denied")
        return forbidden(conf.forbidden_url)
      end
    end

  end
end

return plugin

