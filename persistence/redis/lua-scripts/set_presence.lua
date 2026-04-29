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
local channel = ARGV[5]
local now_ts = tonumber(ARGV[6])
local expiration_ts = tonumber(ARGV[7])

local auth_state = guest and "guest" or "auth"
local expiration_ttl = expiration_ts - now_ts

-- connections registry
redis.call("ZADD", "presence:conns", expiration_ts, conn_id)

-- per connection metadata (username, auth_state, maybe in future device_info or country_code etc)
redis.call("HSET", "presence:conn:" .. conn_id, "user_id", user_id, "username", username, "auth_state", auth_state)
redis.call("EXPIRE", "presence:conn:" .. conn_id, expiration_ttl)

-- user -> connections (which connections belong to which user. user can have multiple ws conns i.e. multiple tabs/devices)
redis.call("SADD", "presence:user:conns:" .. user_id, conn_id)

-- channel -> connections (track which connections are in a channel)
redis.call("ZADD", "presence:channel:conns:" .. channel, expiration_ts, conn_id)
redis.call("EXPIRE", "presence:channel:conns:" .. channel, 86400) -- 24hr

-- conn -> channels (which channels this connection is in)
redis.call("SADD", "presence:conn:channels:" .. conn_id, channel)
