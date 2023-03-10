package dbHelpers

var DefaultSchema = Schema{
	Create: `
    CREATE TABLE IF NOT EXISTS users (
        id text DEFAULT gen_random_uuid() PRIMARY KEY,
        name varchar(255),
        username varchar(255) NOT NULL UNIQUE,
        email varchar(255) NOT NULL UNIQUE,
        password text NOT NULL,
        date_created TIMESTAMP DEFAULT now(),
        date_updated TIMESTAMP DEFAULT now()
    );
`,

	Drop: `
  DROP TABLE IF EXISTS users CASCADE;
  DROP TABLE IF EXISTS expenses CASCADE;
`,
	Alter: `
`,
}