CREATE TABLE `efficient_api`.`messages` (
  `id` SERIAL NOT NULL,
  `title` VARCHAR(100) NULL,
  `body` VARCHAR(250) NULL,
  `created_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
