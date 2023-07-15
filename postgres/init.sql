-- db/modelで定義したテーブル名の変数とかを使って作成したいから
-- goでテーブル作成したいけどsqlxの使い方がわからん
-- created_atとかいるかな？

create table if not exists users(
  discord_id varchar(50) primary key not null,
  name varchar(50) not null,
  is_member boolean not null
);

create table if not exists roles(
  discord_id  varchar(50) not null,
  role_id varchar(50) not null,
  unique(discord_id, role_id)
);

create table if not exists statuses(
  discord_id  varchar(50) not null,
  time_stamp timestamp not null,
  channel_id varchar(50) not null,
  voice_state varchar(50) not null,
  online_status varchar(50) not null,
  unique(discord_id, time_stamp)
);