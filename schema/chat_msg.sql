use foozy_proj;

CREATE TABLE IF NOT EXISTS `foozy_proj`.`chat_msg` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `user` INT NOT NULL,
    `chat` INT NOT NULL,
    `text` TEXT NOT NULL,
    `date` DATETIME NOT NULL,
    `is_read` BOOLEAN NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (user) REFERENCES foozy_proj.auth(id) ON DELETE CASCADE,
    FOREIGN KEY (chat) REFERENCES foozy_proj.chat(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `foozy_proj`.`chat_msg_images` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `parent_msg` INT NOT NULL,
    `path` TEXT NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (parent_msg) REFERENCES chat_msg(id) ON DELETE CASCADE
);
