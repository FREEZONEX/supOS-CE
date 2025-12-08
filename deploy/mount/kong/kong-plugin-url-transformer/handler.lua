local http = require "resty.http"
local ngx = ngx
local cjson = require "cjson"

local plugin = {
    PRIORITY = 1000,  -- 插件的优先级
    VERSION = "1.0.0",  -- 插件版本
}

local auth_url = "http://backend:8080/inter-api/supos/auth/userinfo"

-- 重定向函数
local function authorized(redirect_url)
    if not redirect_url or redirect_url == "" then
        ngx.log(ngx.INFO, "redirect_url 为空，跳过处理")
        return
    else
        ngx.log(ngx.INFO, "重定向到: ", redirect_url)
        return ngx.redirect(redirect_url, ngx.HTTP_MOVED_TEMPORARILY)
    end
end

function plugin:access(conf)
    -- 获取所有 cookies
    local cookies = ngx.req.get_headers()["Cookie"]

    if not cookies then
        ngx.log(ngx.ERR, "未找到 Cookie，拒绝请求")
        return
    end

    -- 提取 supos_community_token 的值
    local supos_community_token = string.match(cookies, "supos_community_token=([^;]+)")
    if not supos_community_token then
        ngx.log(ngx.ERR, "未找到 supos_community_token，拒绝请求")
        return
    end

    -- 向后端发送 GET 请求，带上 supos_community_token
    local httpc = http.new()
    local res, err = httpc:request_uri(auth_url, {
        method = "GET",
        headers = {
            ["Cookie"] = "supos_community_token=" .. supos_community_token,
        },
    })

    if not res then
        ngx.log(ngx.ERR, "请求认证服务失败，错误信息: ", err)
        return
    end

    if res.status ~= 200 then
        ngx.log(ngx.ERR, "认证服务返回非 200 状态码: ", res.status)
        return
    end

    -- 如果认证成功，重定向到 home_url
    ngx.log(ngx.INFO, "认证成功，重定向到主页")
    return authorized(conf.home_url)
end

return plugin