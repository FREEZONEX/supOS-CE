local http = require "resty.http"
local ngx = ngx
local cjson = require "cjson"

local plugin = {
    PRIORITY = 1000,
    VERSION = "1.0.0",
}

local auth_url = "http://uns:8080/inter-api/supos/auth/userinfo"

-- Redirect function
local function authorized(redirect_url)
    if not redirect_url or redirect_url == "" then
        ngx.log(ngx.INFO, "redirect_url is empty, skip processing")
        return
    else
        ngx.log(ngx.INFO, "Redirecting to: ", redirect_url)
        return ngx.redirect(redirect_url, ngx.HTTP_MOVED_TEMPORARILY)
    end
end

function plugin:access(conf)
    -- Get all cookies
    local cookies = ngx.req.get_headers()["Cookie"]

    if not cookies then
        ngx.log(ngx.ERR, "Cookie not found, request denied")
        return
    end

    -- Extract supos_community_token value
    local supos_community_token = string.match(cookies, "supos_community_token=([^;]+)")
    if not supos_community_token then
        ngx.log(ngx.ERR, "supos_community_token not found, request denied")
        return
    end

    -- Send GET request to backend with supos_community_token
    local httpc = http.new()
    local res, err = httpc:request_uri(auth_url, {
        method = "GET",
        headers = {
            ["Cookie"] = "supos_community_token=" .. supos_community_token,
        },
    })

    if not res then
        ngx.log(ngx.ERR, "Request to auth service failed, error: ", err)
        return
    end

    if res.status ~= 200 then
        ngx.log(ngx.ERR, "Auth service returned non-200 status: ", res.status)
        return
    end

    -- If authentication successful, redirect to home_url
    ngx.log(ngx.INFO, "Authentication successful, redirecting to home page")
    return authorized(conf.home_url)
end

return plugin

