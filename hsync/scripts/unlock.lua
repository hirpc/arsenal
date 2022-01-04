local key = KEYS[1]
-- reentrantKey 重入次数
local reentrantKey = KEYS[2]
local reqID = ARGV[1]

if (redis.call("GET", key) == reqID) then
    -- 优先判定是否被重入
    local count = redis.call("GET", reentrantKey)
    if count then
        if (tonumber(count) > 1) then 
            -- 被多次重入，需要逐渐扣除
            count = tonumber(count) - 1
            local live = redis.call("PTTL", reentrantKey)
            redis.call("SET", reentrantKey, count, "PX", live)
            return count
        else
            -- 未被重入，直接删除相关key就好
            redis.call("DEL", reentrantKey)
            redis.call("DEL", key)
            return 0
        end
    else
        redis.call("DEL", key)
        return 0
    end
else
    -- 其他请求，不允许解锁
    return -1
end