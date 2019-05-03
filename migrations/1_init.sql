
-- +migrate Up
CREATE TABLE `histories` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `created` DATETIME NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`)
);

CREATE TABLE `shifts` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `history_id` INTEGER NOT NULL,
    `name` CHAR(10) COMMENT 'キャストの名前',
    `begin` CHAR(5) COMMENT 'hh:mm',
    `end` CHAR(5) COMMENT 'hh:mm',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`history_id`) REFERENCES `histories`(`id`) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE `shifts`;
DROP TABLE `histories`;