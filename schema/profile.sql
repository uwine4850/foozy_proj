CREATE TABLE `foozy_proj`.`auth` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `username` VARCHAR(200) NOT NULL ,
    `password` TEXT NOT NULL ,
    `avatar` TEXT NULL,
    `description` TEXT NULL ,
    `name` VARCHAR(200) NULL ,
    PRIMARY KEY (`id`));