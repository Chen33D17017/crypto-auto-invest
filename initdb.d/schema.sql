use `crypto-db`;

CREATE TABLE IF NOT EXISTS users (
  `uid` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `email` VARCHAR(120) NOT NULL UNIQUE,
  `password` VARCHAR(256) NOT NULL,
  `image_url` VARCHAR(256) NOT NULL DEFAULT '',
  `api_key` VARCHAR(256) NOT NULL DEFAULT '',
  `api_secret` VARCHAR(256) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS wallets(
  `wid` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `type` VARCHAR(10) NOT NULL,
  `amount` FLOAT,
  UNIQUE (`uid`, `type`)
);

CREATE TABLE IF NOT EXISTS trades(
  `tid` VARCHAR(100) NOT NULL PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `from_wid` VARCHAR(100) NOT NULL,
  `from_amount` float NOT NULL,
  `to_wid` VARCHAR(100) NOT NULL,
  `to_amount` float NOT NULL,
  `timestamp` TIMESTAMP NOT NULL
)