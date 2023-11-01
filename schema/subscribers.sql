use foozy_proj;

CREATE TABLE `foozy_proj`.`subscribers` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `subscriber` INT NOT NULL ,
    `profile` INT NOT NULL ,
    PRIMARY KEY (`id`),
    FOREIGN KEY (subscriber) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE,
    FOREIGN KEY (profile) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE)
