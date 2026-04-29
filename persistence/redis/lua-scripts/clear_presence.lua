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

-- get conn channels
local conn_channels = redis.call("SMEMBERS", "presence:conn:channels:" .. conn_id)

-- remove from connections registry
redis.call("ZREM", "presence:conns", conn_id)

-- remove from user set
redis.call("SREM", "presence:user:conns:" .. user_id, conn_id)

-- clean empty user set optionally
if redis.call("SCARD", "presence:user:conns:" .. user_id) == 0 then
  redis.call("DEL", "presence:user:conns:" .. user_id)
end

-- remove from all channel ZSETs
for _, ch in ipairs(conn_channels) do
  redis.call("ZREM", "presence:channel:conns:" .. ch, conn_id)
end

-- delete reverse mapping
redis.call("DEL", "presence:conn:channels:" .. conn_id)

-- delete connection metadata
redis.call("DEL", "presence:conn:" .. conn_id)
