use foozy_proj;

CREATE TABLE IF NOT EXISTS `foozy_proj`.`post_categories` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `name` VARCHAR(200) NOT NULL ,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `foozy_proj`.`posts` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `parent_user` INT NOT NULL ,
    `name` VARCHAR(200) NOT NULL ,
    `description` TEXT NOT NULL ,
    `category` INT NOT NULL ,
    `date` DATETIME NOT NULL ,
    PRIMARY KEY (`id`),
    FOREIGN KEY (parent_user) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE,
    FOREIGN KEY (category) REFERENCES post_categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `foozy_proj`.`post_images` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `postid` INT NOT NULL ,
    `path` TEXT NOT NULL ,
    PRIMARY KEY (`id`),
    FOREIGN KEY (postid) REFERENCES posts(id) ON DELETE CASCADE
);

INSERT IGNORE INTO `foozy_proj`.`post_categories` (`id`, `name`) VALUES ('1', 'Art');
INSERT IGNORE INTO `foozy_proj`.`post_categories` (`id`, `name`) VALUES ('2', 'PixelArt');
INSERT IGNORE INTO `foozy_proj`.`post_categories` (`id`, `name`) VALUES ('3', '3D');
