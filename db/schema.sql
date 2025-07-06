CREATE TYPE "color" AS ENUM ('w', 'b');

CREATE TABLE "user" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game_result" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(20) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game_result_status" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(30) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game_state" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(20) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game_time_category" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(20) NOT NULL,
  "upper_time_limit_secs" integer NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game_time_kind" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(20) NOT NULL,
  "enabled" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game_variant" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(30) NOT NULL,
  "enabled" boolean NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "game" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "white_id" uuid NULL,
  "black_id" uuid NULL,
  "guest_white_id" uuid NULL,
  "guest_black_id" uuid NULL,
  "variant_id" uuid NOT NULL,
  "time_kind_id" uuid NOT NULL,
  "time_category_id" uuid NOT NULL,
  "is_guest" bool NOT NULL,
  "time_control_clock" integer NOT NULL,
  "time_control_increment" integer NOT NULL,
  "reconnect_timeout" integer NOT NULL,
  "first_move_timeout" integer NOT NULL,
  "white_game_clock" integer NOT NULL,
  "black_game_clock" integer NOT NULL,
  "result_id" uuid NULL,
  "result_status_id" uuid NULL,
  "state_id" uuid NOT NULL,
  "start_time" timestamptz NOT NULL,
  "end_time" timestamptz NULL,
  "last_move" timestamptz NULL,
  "fen" character varying(90) NOT NULL,
  "pgn" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "game_black_id_fkey" FOREIGN KEY ("black_id") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_result_fkey" FOREIGN KEY ("result_id") REFERENCES "game_result" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_result_status_fkey" FOREIGN KEY ("result_status_id") REFERENCES "game_result_status" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_state_fkey" FOREIGN KEY ("state_id") REFERENCES "game_state" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_time_category_fkey" FOREIGN KEY ("time_category_id") REFERENCES "game_time_category" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_time_kind_fkey" FOREIGN KEY ("time_kind_id") REFERENCES "game_time_kind" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_variant_fkey" FOREIGN KEY ("variant_id") REFERENCES "game_variant" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "game_white_id_fkey" FOREIGN KEY ("white_id") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CHECK ((white_id IS NOT NULL AND black_id IS NOT NULL AND guest_white_id IS NULL AND guest_black_id IS NULL) OR (guest_white_id IS NOT NULL AND guest_black_id IS NOT NULL AND white_id IS NULL AND black_id IS NULL))
);

CREATE TABLE "game_move" (
  "id" BIGINT GENERATED ALWAYS AS IDENTITY,
  "game_id" uuid NOT NULL,
  "fen" character varying(100) NOT NULL,
  "uci" character varying(5) NOT NULL,
  "san" character varying(10) NOT NULL,
  "played_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "game_move_game_id_fkey" FOREIGN KEY ("game_id") REFERENCES "game" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE TABLE "rating" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "game_time_category_id" uuid NOT NULL, 
  "glicko" integer NOT NULL,
  "glicko2" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  UNIQUE ("user_id", "game_time_category_id"),
  CONSTRAINT "rating_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
