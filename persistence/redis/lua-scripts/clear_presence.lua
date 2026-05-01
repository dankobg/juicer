-- ARGV[1] user_id
-- ARGV[2] conn_id
-- ARGV[3] username
-- ARGV[4] guest
-- ARGV[5] now_ts
-- ARGV[6] expiration_ts

local user_id = ARGV[1]
local conn_id = ARGV[2]
local username = ARGV[3]
local guest = ARGV[4] == "1"
local now_ts = tonumber(ARGV[5])
local expiration_ts = tonumber(ARGV[6])

local auth_state = guest and "guest" or "auth"
local expiration_ttl = expiration_ts - now_ts

local conn_left = {}
local user_left = {}

-- list conn channels
local conn_channels = redis.call("SMEMBERS", "presence:conn:channels:" .. conn_id)

-- remove from connections registry
redis.call("ZREM", "presence:conns", conn_id)

-- remove conn from user connections set
redis.call("SREM", "presence:user:conns:" .. user_id, conn_id)

-- process each channel
for _, ch in ipairs(conn_channels) do
  -- remove conn -> channel
  local conn_channel_removed = redis.call("SREM", "presence:conn:channels:" .. conn_id, ch)
  if conn_channel_removed == 1 then
    table.insert(conn_left, ch)
  end

  -- check if user still has ANY connection in this channel
  local user_conns = redis.call("SMEMBERS", "presence:user:conns:" .. user_id)
  local still_in_channel = false

  for _, other_conn in ipairs(user_conns) do
    if redis.call("SISMEMBER", "presence:conn:channels:" .. other_conn, ch) == 1 then
      still_in_channel = true
      break
    end
  end

  -- if no remaining connections in this channel → remove user
  if not still_in_channel then
    local user_removed = redis.call("SREM", "presence:user:channels:" .. user_id, ch)
    if user_removed == 1 then
      table.insert(user_left, ch)

      redis.call("SREM", "presence:channel:users:" .. ch, user_id)
    end
  end
end

-- clean empty user set optionally
if redis.call("SCARD", "presence:user:conns:" .. user_id) == 0 then
  redis.call("DEL", "presence:user:conns:" .. user_id)
end

-- delete reverse mapping
redis.call("DEL", "presence:conn:channels:" .. conn_id)

-- delete connection metadata
redis.call("DEL", "presence:conn:meta:" .. conn_id)

-- delete user metadata
redis.call("DEL", "presence:user:meta:" .. user_id)

-- update user last seen
if not guest then
  redis.call("ZADD", "presence:user:last-seen", now_ts, user_id)
end

return { conn_left, user_left }
