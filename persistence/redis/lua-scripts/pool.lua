-- ARGV[1] pool_key

local pool_key = KEYS[1]

local player_ids = redis.call('ZRANGE', pool_key, 0, 1)

if #player_ids == 2 then
  redis.call('ZREM', pool_key, player_ids[1], player_ids[2])
  return player_ids
end

return {}
