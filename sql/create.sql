create table tb_user_song_tags (
  pk_user_song_tag_id integer primary key,
  fk_user_id INTEGER not null,
  fk_song_id text not null,
  fk_tag_id integer not null
);


create table tb_tag (
  pk_tag_id integer primary key,
  tag_name text not null
);
