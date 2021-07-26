CREATE TABLE `chats` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `sender` varchar(100) NOT NULL,
    `receiver` varchar(100) NOT NULL,
    `body` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4