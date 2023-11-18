use foozy_proj;

CREATE TABLE IF NOT EXISTS `foozy_proj`.`chat` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `user1` INT NOT NULL,
    `user2` INT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (user1) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE,
    FOREIGN KEY (user2) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE);
