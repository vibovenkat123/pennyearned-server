CREATE TABLE IF NOT EXISTS expenses (
    id varchar(36) PRIMARY KEY NOT NULL,
    owner_id varchar(36) NOT NULL,
    name varchar(255) NOT NULL,
    spent int(11) DEFAULT 0 NOT NULL,
    date_created DATETIME DEFAULT NOW(),
    date_updated DATETIME DEFAULT NOW()
);
