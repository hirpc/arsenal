local key = KEYS[1]
local reqID = ARGV[1]
local expireMS = ARGV[2]

local res = redis.call("SET", key, reqID, "PX", expireMS, "NX")
if res == false then 
    return 0
elseif res["ok"] == "OK" then
    return 1
elseif redis.call("TTL", key) == -1 then
    redis.call("PEXPIRE", key, expireMS)
end
return 0