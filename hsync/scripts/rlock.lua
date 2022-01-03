local readKey = "read_" .. KEYS[1]
local writeKey = "write_" .. KEYS[2]
local reqID = ARGV[1]
local expireMS = ARGV[2]

-- hash表写入
local function HSET(key, field, value)
    return redis.call("HSET", key, field, value)
end
-- TTL
local function PTTL(key)
    return redis.call("PTTL", key)
end
-- PEXPIRE
local function PEXPIRE(key, t)
    return redis.call("PEXPIRE", key, t)
end
-- UPDATE_EXPIRATION 更新某个锁的过期时间
-- 尝试与传入的过期时间取最大值
local function UPDATE_EXPIRATION(key)
    PEXPIRE(key, math.max(PTTL(key),expireMS))
end

if not redis.call("GET", writeKey) then
    -- 没有写锁，可以加读锁
    local count = redis.call("HGET", readKey, reqID)
    if count then
        count = tonumber(count) + 1
        HSET(readKey, reqID, count)
    else
        HSET(readKey, reqID, 1)
    end
    -- 更新key的过期时间
    UPDATE_EXPIRATION(readKey)
    return 1
else
    -- 存在写锁，不允许加读锁
    return 0
end