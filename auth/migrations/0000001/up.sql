CREATE TABLE IF NOT EXISTS users (
    id varchar(36) PRIMARY KEY NOT NULL,
    username varchar(255) NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    date_created DATETIME DEFAULT NOW(),
    date_updated DATETIME DEFAULT NOW()
);
