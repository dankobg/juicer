-- ARGV[1] pool_key

local pool_key = KEYS[1]

local user_ids = redis.call('ZRANGE', pool_key, 0, 1)

if #user_ids == 2 then
  redis.call("ZREM", pool_key, user_ids[1], user_ids[2])

  for _, user_id in ipairs(user_ids) do
    local pool_user_key = "pool-user:" .. user_id
    local pool_user_meta_key = "pool-user:meta:" .. user_id

    redis.call("DEL", pool_user_key)
    redis.call("DEL", pool_user_meta_key)
  end

  return user_ids
end

return {}
