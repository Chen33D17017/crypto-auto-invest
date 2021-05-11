use `crypto-db`;
-- create table users
-- (
-- id int auto_increment primary key,
-- username varchar(64) unique not null,
-- email varchar(120) unique not null,
-- password_hash varchar(128) not null,
-- avatar varchar(128) not null
-- );
-- insert into user values(1, "zhangsan","test12345@qq.com","passwd","avaterpath");
-- insert into user values(2, "lisi","12345test@qq.com","passwd","avaterpath");

CREATE TABLE IF NOT EXISTS users (
  `uid` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `email` VARCHAR(120) NOT NULL UNIQUE,
  `password` VARCHAR(256) NOT NULL,
  `image_url` VARCHAR(256) NOT NULL
);