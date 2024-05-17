CREATE TABLE `courses` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT, 
    
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NULL DEFAULT '',
    
    `description` varchar(128) COLLATE utf8mb4_unicode_ci NULL DEFAULT '',
    
    `price` float NULL DEFAULT NULL,

    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT NULL,
    
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;