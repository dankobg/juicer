-- +goose Up
-- +goose StatementBegin
insert into "game_variant" ("name", "enabled") values 
  ('standard', true),
  ('atomic', false),
  ('crazyhouse', false),
  ('chess960', false),
  ('king-of-the-hill', false),
  ('three-check', false),
  ('horde', false),
  ('racing-kings', false);

insert into "game_time_kind" ("name", "enabled") values 
  ('realtime', true),
  ('correspondence', false),
  ('unlimited', false);

insert into "game_time_category" ("name", "upper_time_limit_secs") values 
  ('hyperbullet', 60),
  ('bullet', 180),
  ('blitz', 600),
  ('rapid', 1800),
  ('classical', null);

insert into "game_result" ("name") values 
  ('white-won'),
  ('black-won'),
  ('draw'),
  ('interrupted');

insert into "game_result_status" ("name") values 
  ('checkmate'), 
  ('insufficient-material'), 
  ('threefold-repetition'), 
  ('fivefold-repetition'), 
  ('fifty-move-rule'), 
  ('seventyfive-move-rule'), 
  ('stalemate'), 
  ('resignation'), 
  ('draw-agreed'), 
  ('flagged'), 
  ('adjudication'), 
  ('timed-out'), 
  ('aborted'), 
  ('interrupted');

insert into "game_state" ("name") values 
  ('idle'), 
  ('waiting-start'), 
  ('in-progress'), 
  ('finished'), 
  ('interrupted');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table "game_variant" cascacde;
truncate table "game_time_kind" cascacde;
truncate table "game_time_category" cascacde;
truncate table "game_result" cascacde;
truncate table "game_result_status" cascacde;
truncate table "game_state" cascacde;
-- +goose StatementEnd