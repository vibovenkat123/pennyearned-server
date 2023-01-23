package dbHelpers

var defaultSchema = Schema{
	create: `
CREATE TABLE IF NOT EXISTS users (
    id text DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS expenses (
    id text DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id text REFERENCES users(id) NOT NULL,
    name varchar(255) NOT NULL,
    spent integer DEFAULT 0 NOT NULL
);`,

	drop: `
  DROP TABLE IF EXISTS users CASCADE;
  DROP TABLE IF EXISTS expenses CASCADE;
`,
	alter: `
`,
}
