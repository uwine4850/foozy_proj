use foozy_proj;

CREATE TABLE IF NOT EXISTS `foozy_proj`.`chat_msg_count` (
    `id` INT NOT NULL AUTO_INCREMENT ,
    `chat` INT NOT NULL ,
    `user` INT NOT NULL ,
    `count` INT NULL DEFAULT '0' ,
    PRIMARY KEY (`id`),
    FOREIGN KEY (chat) REFERENCES foozy_proj.chat(id) ON DELETE CASCADE,
    FOREIGN KEY (user) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE
);