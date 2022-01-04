local readKey = "read_" .. KEYS[1]
local reqID = ARGV[1]

local function HSET(key, field, value)
    return redis.call("HSET", key, field, value)
end

local function HGET(key, field)
    return redis.call("HGET", key, field)
end

local function HDEL(key, field)
    return redis.call("HDEL", key, field)
end

local count = HGET(readKey, reqID)
if count then
    if (tonumber(count) > 1) then
        count = tonumber(count) - 1
        HSET(readKey, reqID, count)
    else
        HDEL(readKey, reqID)
    end
end