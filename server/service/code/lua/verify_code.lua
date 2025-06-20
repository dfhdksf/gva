local key = KEYS[1]
-- 使用次数，也就是验证次数
local cntKey = key .. ":cnt"
-- 预期中的验证码：用户输入的验证码
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)
-- 验证次数已经耗尽了
if cnt <= 0 then
    return -1
end
-- 验证码相等
if code == expectedCode then
    -- 把次数标记位 -1，认为验证码不可用
    redis.call("set", cntKey, -1)
    return 0
else
    -- 可能使用户手一抖输错了
    redis.call("decr", cntKey)
    return -2
end