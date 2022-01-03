local readKey = "read_" .. KEYS[1]
local reqID = ARGV[1]

local count = redis.call("HGET", readKey, reqID)
if count then
    if (tonumber(count) > 1) then
        count = tonumber(count) - 1
        redis.call("HSET", readKey, reqID, count)
    else
        redis.call("HDEL", readKey, reqID)
    end
end