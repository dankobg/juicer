-- ARGV[1] user_id
-- ARGV[2] conn_id
-- ARGV[3] username
-- ARGV[4] guest
-- ARGV[5] now_ts
-- ARGV[6] expiration_ts

local user_id = ARGV[1]
local conn_id = ARGV[2]
local username = ARGV[3]
local guest = ARGV[4] == "true"
local now_ts = tonumber(ARGV[5])
local expiration_ts = tonumber(ARGV[6])

local auth_state = guest and "guest" or "auth"
local expiration_ttl = expiration_ts - now_ts

-- user last seen
redis.call("ZADD", "presence:user:last-seen", now_ts, user_id)

-- refresh connection
redis.call("ZADD", "presence:conns", expiration_ts, conn_id)

-- refresh expiry
redis.call("EXPIRE", "presence:conn:" .. conn_id, expiration_ttl)
