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

CREATE TABLE IF NOT EXISTS crypto_name (
  `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(10) NOT NULL
);

INSERT INTO `crypto_name`(`name`) VALUES("jpy");
INSERT INTO `crypto_name`(`name`) VALUES("eth");
INSERT INTO `crypto_name`(`name`) VALUES("btc");
INSERT INTO `crypto_name`(`name`) VALUES("xrp");
INSERT INTO `crypto_name`(`name`) VALUES("ltc");
INSERT INTO `crypto_name`(`name`) VALUES("mona");
INSERT INTO `crypto_name`(`name`) VALUES("bcc");
INSERT INTO `crypto_name`(`name`) VALUES("xlm");
INSERT INTO `crypto_name`(`name`) VALUES("qtum");
INSERT INTO `crypto_name`(`name`) VALUES("bat");


CREATE TABLE IF NOT EXISTS wallets(
  `wid` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `crypto_id` INT NOT NULL,
  `amount` FLOAT NOT NULL,

  FOREIGN KEY (`uid`) REFERENCES `users` (`uid`),
  FOREIGN KEY (`crypto_id`) REFERENCES `crypto_name`(`id`),
  UNIQUE (`uid`, `crypto_id`)
);



CREATE TABLE IF NOT EXISTS orders(
  `oid` VARCHAR(100) NOT NULL PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `from_wid` VARCHAR(100) NOT NULL,
  `from_amount` float NOT NULL,
  `to_wid` VARCHAR(100) NOT NULL,
  `to_amount` float NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  `fee` float NOT NULL,
  `type` VARCHAR(10),

  FOREIGN KEY (`from_wid`) REFERENCES `wallets` (`wid`),
  FOREIGN KEY (`to_wid`) REFERENCES `wallets` (`wid`)
);

CREATE TABLE IF NOT EXISTS crons(
  `id` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY, 
  `uid` VARCHAR(100) NOT NULL,
  `crypto_id` INT NOT NULL,
  `amount` FLOAT NOT NULL,
  `time_pattern` VARCHAR(100) NOT NULL,

  FOREIGN KEY (`uid`) REFERENCES `users` (`uid`),
  FOREIGN KEY (`crypto_id`) REFERENCES `crypto_name`(`id`),
  UNIQUE (`uid`, `crypto_id`, `time_pattern`)
);

CREATE TABLE IF NOT EXISTS auto_trades(
  `id` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `crypto_id` INT NOT NULL,

  FOREIGN KEY (`uid`) REFERENCES `users` (`uid`),
  FOREIGN KEY (`crypto_id`) REFERENCES `crypto_name`(`id`)
);

-- view for wallets
CREATE VIEW `wallets_view` AS
SELECT
    w.wid as `wid`,
    w.uid as `uid`,
    ct.name as `crypto_name`,
    w.amount as `amount`
FROM
    wallets w
INNER JOIN
    crypto_name ct
ON
    w.crypto_id=ct.id;


-- view for corns
CREATE VIEW `crons_view` AS
SELECT
    c.id as id,
    c.uid as `uid`,
    ct.name as `crypto_name`,
    c.amount as amount,
    c.time_pattern as time_pattern
FROM
    crons c
INNER JOIN
    crypto_name ct
ON
    c.crypto_id=ct.id;

-- view of auto_trades
CREATE VIEW `auto_trades_view` AS
SELECT
    `at`.id as id,
    `at`.`uid` as `uid`,
    ct.name as `crypto_name`
FROM
    auto_trades as `at`
INNER JOIN
    crypto_name ct
ON
    `at`.crypto_id=ct.id;