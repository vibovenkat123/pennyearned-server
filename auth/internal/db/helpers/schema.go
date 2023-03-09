package dbHelpers

var defaultSchema = Schema{
	create: `
    CREATE TABLE IF NOT EXISTS users (
        id text DEFAULT gen_random_uuid() PRIMARY KEY,
        name varchar(255) NOT NULL,
        username varchar(255) NOT NULL UNIQUE,
        email varchar(255) NOT NULL UNIQUE,
        password text NOT NULL,
        date_created TIMESTAMP DEFAULT now(),
        date_updated TIMESTAMP DEFAULT now()
    );
`,

	drop: `
  DROP TABLE IF EXISTS users CASCADE;
  DROP TABLE IF EXISTS expenses CASCADE;
`,
	alter: `
`,
}