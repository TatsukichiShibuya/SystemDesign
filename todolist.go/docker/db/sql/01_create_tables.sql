-- Table for tasks
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `owners`;

CREATE TABLE `tasks` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `title` varchar(50) NOT NULL,
    `is_done` boolean NOT NULL DEFAULT b'0',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `deadline` datetime NOT NULL  DEFAULT CURRENT_TIMESTAMP,
    `detail` varchar(1000) NOT NULL DEFAULT "",
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `username` varchar(20) NOT NULL,
    `passward` varchar(256) NOT NULL,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `owners` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `username` varchar(20) NOT NULL,
    `taskid` bigint(20) NOT NULL,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;
