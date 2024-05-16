CREATE TABLE IF NOT EXISTS `helps` (
    id int AUTO_INCREMENT NOT NULL,
    title varchar(100) NOT NULL,
    content text NOT NULL,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    PRIMARY KEY(id)
);