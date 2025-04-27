-- 发送到的 key，也就是 code:业务:手机号码
local key = KEYS[1]
-- 使用次数，也就是验证次数
local cntKey = key..":cnt"
-- val为验证码
local val = ARGV[1]
-- 验证码的有效时间(秒)
local expire = tonumber(ARGV[2])
-- 重发间隔(秒)
local resendInterval = tonumber(ARGV[3])
-- 调用ttl命令，返回-1/-2等，tonumber是将结果转为数字
local ttl = tonumber(redis.call("ttl", key))

-- -1 是 key 存在，但是没有过期时间
if ttl == -1 then
    -- 有人误操作，导致 key 冲突
    return -2 -- 表示系统错误
-- -2 是 key 不存在，ttl < (expire - resendInterval) 是发了一个验证码，已经超过重发间隔了，可以重新发送
elseif ttl == -2 or ttl < (expire - resendInterval) then
    redis.call("set", key, val)
    redis.call("expire", key, expire)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, expire)
    -- 完美，符合预期
    return 0
else
    -- 已经发送了一个验证码，但是还不到重发间隔
    return -1
end