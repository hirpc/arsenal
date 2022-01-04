local key = KEYS[1]
-- reentrantKey 重入次数
local reentrantKey = KEYS[2]
local reqID = ARGV[1]
local expireMS = ARGV[2]

if redis.call("SET", key, reqID, "PX", expireMS, "NX") then
    redis.call("SET", reentrantKey, 1, "PX", expireMS)
    return 1
else
    -- 存在写锁
    if (redis.call("GET", key) == reqID) then 
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
