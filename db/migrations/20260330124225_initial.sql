-- +goose up
-- +goose StatementBegin
create type "color" as enum ('w', 'b');

create table "user" (
  "id" uuid not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_user_id" primary key ("id")
);

create table "friendship" (
  "initiator_id" uuid not null,
  "receiver_id" uuid not null,
  "status" varchar(30) not null default 'pending',
  "created_at" timestamptz not null default current_timestamp,
  "answered_at" timestamptz,
  constraint "pk_friendship_initiator_id_receiver_id" primary key ("initiator_id", "receiver_id"),
  constraint "ck_friendship_no_self" check ("initiator_id" <> "receiver_id"),
  constraint "ck_friendship_status_valid" check ("status" in ('pending', 'accepted', 'declined')),
  constraint "fk_friendship_initiator_id" foreign key ("initiator_id") references "user" ("id"),
  constraint "fk_friendship_receiver_id" foreign key ("receiver_id") references "user" ("id")
);

create unique index "uq_friendship_users_undirected" on "friendship" (
  least("initiator_id", "receiver_id"),
  greatest("initiator_id", "receiver_id")
);

create table "following" (
  "user_id" uuid not null,
  "followed_user_id" uuid not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_following_user_id_followed_user_id" primary key ("user_id", "followed_user_id"),
  constraint "ck_following_no_self" check ("user_id" <> "followed_user_id"),
  constraint "uq_following_user_id_followed_user_id" unique ("user_id", "followed_user_id"),
  constraint "fk_following_user_id" foreign key ("user_id") references "user" ("id"),
  constraint "fk_following_followed_user_id" foreign key ("followed_user_id") references "user" ("id")
);

create table "blocklist" (
  "user_id" uuid not null,
  "blocked_user_id" uuid not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_blocklist_user_id_blocked_user_id" primary key ("user_id", "blocked_user_id"),
  constraint "ck_blocklist_no_self" check ("user_id" <> "blocked_user_id"),
  constraint "uq_blocklist_user_id_blocked_user_id" unique ("user_id", "blocked_user_id"),
  constraint "fk_blocklist_user_id" foreign key ("user_id") references "user" ("id"),
  constraint "fk_blocklist_blocked_user_id" foreign key ("blocked_user_id") references "user" ("id")
);

create table "game_result" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_result_id" primary key ("id")
);

create table "game_result_status" (
  "id" bigint not null generated always as identity,
  "name" character varying(30) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_result_status_id" primary key ("id")
);

create table "game_state" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_state_id" primary key ("id")
);

create table "game_time_category" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "upper_time_limit_secs" integer null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_time_category_id" primary key ("id")
);

create table "game_time_kind" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "enabled" boolean not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_time_kind_id" primary key ("id")
);

create table "game_variant" (
  "id" bigint not null generated always as identity,
  "name" character varying(30) not null,
  "enabled" boolean not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_variant_id" primary key ("id")
);

create table "game" (
  "id" bigint not null generated always as identity,
  "white_id" uuid,
  "black_id" uuid,
  "white_is_guest" bool not null,
  "black_is_guest" bool not null,
  "guest_white_id" uuid,
  "guest_black_id" uuid,
  "game_variant_id" bigint not null,
  "game_time_kind_id" bigint not null,
  "game_time_category_id" bigint not null,
  "game_state_id" bigint not null,
  "game_result_id" bigint,
  "game_result_status_id" bigint,
  "time_control_clock_ms" integer not null,
  "time_control_increment_ms" integer not null,
  "reconnect_timeout_ms" integer not null,
  "first_move_timeout_ms" integer not null,
  "white_game_clock" integer not null,
  "black_game_clock" integer not null,
  "rated" boolean not null,
  "start_time" timestamptz,
  "end_time" timestamptz,
  "last_move" timestamptz,
  "fen" character varying(90) not null,
  "pgn" text,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_game_id" primary key ("id"),
  constraint "fk_game_white_id" foreign key ("white_id") references "user" ("id"),
  constraint "fk_game_black_id" foreign key ("black_id") references "user" ("id"),
  constraint "fk_game_variant_id" foreign key ("game_variant_id") references "game_variant" ("id"),
  constraint "fk_game_time_kind_id" foreign key ("game_time_kind_id") references "game_time_kind" ("id"),
  constraint "fk_game_time_category_id" foreign key ("game_time_category_id") references "game_time_category" ("id"),
  constraint "fk_game_state_id" foreign key ("game_state_id") references "game_state" ("id"),
  constraint "fk_game_result_id" foreign key ("game_result_id") references "game_result" ("id"),
  constraint "fk_game_result_status_id" foreign key ("game_result_status_id") references "game_result_status" ("id"),
  constraint "ck_game_white_black" check (
    ("white_id" is not null and "black_id" is not null and "guest_white_id" is null and "guest_black_id" is null) or
    ("guest_white_id" is not null and "guest_black_id" is not null and "white_id" is null and "black_id" is null)
  )
);

create table "game_move" (
  "id" bigint generated always as identity,
  "game_id" bigint not null,
  "fen" character varying(100) not null,
  "uci" character varying(5) not null,
  "san" character varying(10) not null,
  "played_at" timestamptz not null default current_timestamp,
  constraint "pk_game_move_id" primary key ("id"),
  constraint "fk_game_move_game_id" foreign key ("game_id") references "game" ("id") on delete cascade
);

create table "rating" (
  "id" bigint not null generated always as identity,
  "user_id" uuid not null,
  "game_time_category_id" bigint not null,
  "glicko" integer not null,
  "glicko2" integer not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  constraint "pk_rating_id" primary key ("id"),
  constraint "uq_rating_user_id_game_time_category_id" unique ("user_id", "game_time_category_id"),
  constraint "fk_rating_user_id" foreign key ("user_id") references "user" ("id") on update cascade on delete cascade
);
-- +goose StatementEnd

-- +goose down
-- +goose StatementBegin
drop table if exists "user" cascade;
drop table if exists "game_result" cascade;
drop table if exists "game_result_status" cascade;
drop table if exists "game_state" cascade;
drop table if exists "game_time_category" cascade;
drop table if exists "game_time_kind" cascade;
drop table if exists "game_variant" cascade;
drop table if exists "game" cascade;
drop table if exists "game_move" cascade;
drop table if exists "rating" cascade;
drop type if exists "color";
-- +goose StatementEnd
