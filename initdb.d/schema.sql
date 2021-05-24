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

CREATE TABLE IF NOT EXISTS currency_type (
  `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(10) NOT NULL
);

INSERT INTO `currency_type`(`name`) VALUES("jpy");
INSERT INTO `currency_type`(`name`) VALUES("eth");
INSERT INTO `currency_type`(`name`) VALUES("btc");
INSERT INTO `currency_type`(`name`) VALUES("xrp");
INSERT INTO `currency_type`(`name`) VALUES("ltc");
INSERT INTO `currency_type`(`name`) VALUES("mona");
INSERT INTO `currency_type`(`name`) VALUES("bcc");
INSERT INTO `currency_type`(`name`) VALUES("xlm");
INSERT INTO `currency_type`(`name`) VALUES("qtum");
INSERT INTO `currency_type`(`name`) VALUES("bat");


CREATE TABLE IF NOT EXISTS wallets(
  `wid` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `type_id` INT NOT NULL,
  `amount` FLOAT NOT NULL,

  FOREIGN KEY (`uid`) REFERENCES `users` (`uid`),
  FOREIGN KEY (`type_id`) REFERENCES `currency_type`(`id`),
  UNIQUE (`uid`, `type_id`)
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
  `type_id` INT NOT NULL,
  `amount` FLOAT NOT NULL,
  `time_pattern` VARCHAR(100) NOT NULL,

  FOREIGN KEY (`uid`) REFERENCES `users` (`uid`),
  FOREIGN KEY (`type_id`) REFERENCES `currency_type`(`id`),
  UNIQUE (`uid`, `type_id`, `time_pattern`)
);

CREATE TABLE IF NOT EXISTS auto_trades(
  `id` VARCHAR(100) DEFAULT (uuid()) PRIMARY KEY,
  `uid` VARCHAR(100) NOT NULL,
  `type_id` INT NOT NULL,

  FOREIGN KEY (`type_id`) REFERENCES `currency_type`(`id`),
  UNIQUE (`uid`, `type_id`)
);

-- view for wallets
CREATE VIEW `wallets_view` AS
SELECT
    w.wid as `wid`,
    w.uid as `uid`,
    ct.name as `type`,
    w.amount as `amount`
FROM
    wallets w
INNER JOIN
    currency_type ct
ON
    w.type_id=ct.id;


-- view for corns
CREATE VIEW `crons_view` AS
SELECT
    c.id as id,
    c.uid as `uid`,
    ct.name as `type`,
    c.amount as amount,
    c.time_pattern as time_pattern
FROM
    crons c
INNER JOIN
    currency_type ct
ON
    c.type_id=ct.id;

-- view of auto_trades
CREATE VIEW `auto_trades_view` AS
SELECT
    `at`.id as id,
    `at`.`uid` as `uid`,
    ct.name as `type`
FROM
    auto_trades as `at`
INNER JOIN
    currency_type ct
ON
    `at`.type_id=ct.id;