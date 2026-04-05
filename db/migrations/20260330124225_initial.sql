-- +goose up
-- +goose StatementBegin
create type "color" as enum ('w', 'b');

create table "user" (
  "id" uuid not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game_result" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game_result_status" (
  "id" bigint not null generated always as identity,
  "name" character varying(30) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game_state" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game_time_category" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "upper_time_limit_secs" integer null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game_time_kind" (
  "id" bigint not null generated always as identity,
  "name" character varying(20) not null,
  "enabled" boolean not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game_variant" (
  "id" bigint not null generated always as identity,
  "name" character varying(30) not null,
  "enabled" boolean not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id")
);

create table "game" (
  "id" bigint not null generated always as identity,
  "white_id" uuid null,
  "black_id" uuid null,
  "guest_white_id" uuid null,
  "guest_black_id" uuid null,
  "variant_id" bigint not null,
  "time_kind_id" bigint not null,
  "time_category_id" bigint not null,
  "is_guest" bool not null,
  "time_control_clock" integer not null,
  "time_control_increment" integer not null,
  "reconnect_timeout" integer not null,
  "first_move_timeout" integer not null,
  "white_game_clock" integer not null,
  "black_game_clock" integer not null,
  "result_id" bigint null,
  "result_status_id" bigint null,
  "state_id" bigint not null,
  "start_time" timestamptz not null,
  "end_time" timestamptz null,
  "last_move" timestamptz null,
  "fen" character varying(90) not null,
  "pgn" text null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "game_black_id_fkey" foreign key ("black_id") references "user" ("id") on update no action on delete no action,
  constraint "game_result_fkey" foreign key ("result_id") references "game_result" ("id") on update no action on delete no action,
  constraint "game_result_status_fkey" foreign key ("result_status_id") references "game_result_status" ("id") on update no action on delete no action,
  constraint "game_state_fkey" foreign key ("state_id") references "game_state" ("id") on update no action on delete no action,
  constraint "game_time_category_fkey" foreign key ("time_category_id") references "game_time_category" ("id") on update no action on delete no action,
  constraint "game_time_kind_fkey" foreign key ("time_kind_id") references "game_time_kind" ("id") on update no action on delete no action,
  constraint "game_variant_fkey" foreign key ("variant_id") references "game_variant" ("id") on update no action on delete no action,
  constraint "game_white_id_fkey" foreign key ("white_id") references "user" ("id") on update no action on delete no action,
  check (("white_id" is not null and "black_id" is not null and "guest_white_id" is null and "guest_black_id" is null) or ("guest_white_id" is not null and "guest_black_id" is not null and "white_id" is null and "black_id" is null))
);

create table "game_move" (
  "id" bigint generated always as identity,
  "game_id" bigint not null,
  "fen" character varying(100) not null,
  "uci" character varying(5) not null,
  "san" character varying(10) not null,
  "played_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  constraint "game_move_game_id_fkey" foreign key ("game_id") references "game" ("id") on update no action on delete cascade
);

create table "rating" (
  "id" bigint not null generated always as identity,
  "user_id" uuid not null,
  "game_time_category_id" bigint not null, 
  "glicko" integer not null,
  "glicko2" integer not null,
  "created_at" timestamptz not null default current_timestamp,
  "updated_at" timestamptz not null default current_timestamp,
  primary key ("id"),
  unique ("user_id", "game_time_category_id"),
  constraint "rating_user_id_fkey" foreign key ("user_id") references "user" ("id") on update cascade on delete cascade
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
