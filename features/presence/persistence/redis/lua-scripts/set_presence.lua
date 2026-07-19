-- ARGV[1] user_id
-- ARGV[2] conn_id
-- ARGV[3] username
-- ARGV[4] guest
-- ARGV[5] channel
-- ARGV[6] now_ts
-- ARGV[7] expiration_ts

local user_id = ARGV[1]
local conn_id = ARGV[2]
local username = ARGV[3]
local guest = ARGV[4] == "1"
local now_ts = tonumber(ARGV[5])
local expiration_ts = tonumber(ARGV[6])
local channels_length = tonumber(ARGV[7])

local channels = {}

for i = 1, channels_length do
  channels[i] = ARGV[7 + i]
end

local auth_state = guest and "guest" or "auth"
local expiration_ttl = expiration_ts - now_ts

local conn_joined = {}
local user_joined = {}

-- track connections
redis.call("ZADD", "presence:conns", expiration_ts, conn_id)

-- connection metadata (username, auth_state, maybe in future device_info or country_code etc)
redis.call("HSET", "presence:conn:meta:" .. conn_id,
  "user_id", user_id,
  "username", username,
  "auth_state", auth_state
)
redis.call("EXPIRE", "presence:conn:meta:" .. conn_id, expiration_ttl)

-- user metadata
redis.call("HSET", "presence:user:meta:" .. user_id,
  "username", username,
  "auth_state", auth_state
)
redis.call("EXPIRE", "presence:user:meta:" .. user_id, expiration_ttl)

-- user -> connections (which connections belong to which user. user can have multiple ws conns i.e. multiple tabs/devices)
redis.call("SADD", "presence:user:conns:" .. user_id, conn_id)

for _, channel in ipairs(channels) do
  -- conn -> channels (which channels this connection is in)
  local conn_channel_added = redis.call("SADD", "presence:conn:channels:" .. conn_id, channel)
  if conn_channel_added == 1 then
    table.insert(conn_joined, channel)
  end

  -- user -> channels (which channels this user is in)
  local user_channel_added = redis.call("SADD", "presence:user:channels:" .. user_id, channel)
  if user_channel_added == 1 then
    table.insert(user_joined, channel)
  end

  -- channel -> users (track which users are in a channel)
  if user_channel_added == 1 then
    redis.call("SADD", "presence:channel:users:" .. channel, user_id)
    redis.call("EXPIRE", "presence:channel:users:" .. channel, 86400)
  end
end

return { conn_joined, user_joined }
