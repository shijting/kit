local val = redis.call('get', KEYS[1])
if val == false then
    return redis.call('set', KEYS[1], ARGV[1], 'PX', ARGV[2])
elseif val == ARGV[1] then
    redis.call('expire', KEYS[1], ARGV[2])
    return "OK"
else
    return ""
end