ALTER TABLE orders
ADD COLUMN `strategy` INT NOT NULL DEFAULT 0;

ALTER TABLE orders
DROP COLUMN `type`;