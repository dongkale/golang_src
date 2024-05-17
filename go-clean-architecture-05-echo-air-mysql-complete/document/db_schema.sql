CREATE TABLE `article` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT, 
    
    `title` varchar(255) COLLATE utf8mb4_unicode_ci NULL DEFAULT '',
    
    `content` varchar(128) COLLATE utf8mb4_unicode_ci NULL DEFAULT '',
    
    `author_id` bigint NULL DEFAULT NULL,

    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT NULL,
    
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO test_02.article
(title, content, author_id, created_at, updated_at)
VALUES('title_01', 'content_01', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

select * from article
select * from author

CREATE TABLE `author` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT, 
    
    `name` varchar(255) COLLATE utf8mb4_unicode_ci NULL DEFAULT '',    
    
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT NULL,
    
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO test_02.author
(name, created_at, updated_at)
VALUES('name_01', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);