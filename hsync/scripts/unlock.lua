local key = KEYS[1]
local reqID = ARGV[1]

if redis.call("GET", key) == reqID then
    redis.call("DEL", key)
    return 1
end
return 0