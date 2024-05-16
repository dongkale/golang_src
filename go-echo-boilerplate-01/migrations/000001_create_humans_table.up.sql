CREATE TABLE IF NOT EXISTS `humans` (
    id int AUTO_INCREMENT NOT NULL,
    name varchar(100) NOT NULL,
    race varchar(100) NOT NULL,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW(),
    PRIMARY KEY(id)
);