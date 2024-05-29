CREATE TABLE IF NOT EXISTS `users` (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `first_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `pass_hash` BLOB NOT NULL,
    `verification_code` VARCHAR(255) DEFAULT NULL,
    `verified` BOOLEAN DEFAULT false,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;