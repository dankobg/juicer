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

-- refresh connection
redis.call("ZADD", "presence:conns", expiration_ts, conn_id)

-- refresh conn metadata expiration
redis.call("EXPIRE", "presence:conn:meta:" .. conn_id, expiration_ttl)

-- refresh user metadata expiration
redis.call("EXPIRE", "presence:user:meta:" .. user_id, expiration_ttl)

-- refresh conn -> channels expiration
redis.call("EXPIRE", "presence:conn:channels:" .. conn_id, expiration_ttl)

-- list channels to refresh channel:user TTLs
local conn_channels = redis.call("SMEMBERS", "presence:conn:channels:" .. conn_id)

for _, ch in ipairs(conn_channels) do
  -- refresh channel users TTL (same as set_presence: 86400)
  redis.call("EXPIRE", "presence:channel:users:" .. ch, 86400)
end

-- user last seen
if not guest then
  redis.call("ZADD", "presence:user:last-seen", now_ts, user_id)
end
