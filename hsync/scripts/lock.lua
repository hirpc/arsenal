local writeKey = KEYS[1]
local readKey = KEYS[2]
-- reentrantKey 重入次数
local reentrantKey = KEYS[3]
local reqID = ARGV[1]
local expireMS = ARGV[2]

-- 首先检查是否存在读锁，如果存在，则不允许加写锁
if redis.call("HGET", readKey, reqID) then
    return 0
end
-- 无读锁，尝试写锁
if redis.call("SET", writeKey, reqID, "PX", expireMS, "NX") then
    redis.call("SET", reentrantKey, 1, "PX", expireMS)
    return 1
else
    -- 存在写锁
    if (redis.call("GET", writeKey) == reqID) then 
        -- 当前请求，重入锁
        local count = redis.call("GET", reentrantKey)
        if count then
            count = tonumber(count) + 1
            redis.call("SET", reentrantKey, count, "PX", expireMS)
            return count
        else
            -- 正常应该不会走到这里，兜底
            redis.call("SET", reentrantKey, 1, "PX", expireMS)
            return 1
        end
    else
        return 0
    end
end
