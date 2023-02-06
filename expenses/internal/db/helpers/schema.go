package dbHelpers

var defaultSchema = Schema{
	create: `
CREATE TABLE IF NOT EXISTS expenses (
    id text DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id text NOT NULL,
    name varchar(255) NOT NULL,
    spent integer DEFAULT 0 NOT NULL,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP DEFAULT now()
);`,

	drop: `
  DROP TABLE IF EXISTS users CASCADE;
  DROP TABLE IF EXISTS expenses CASCADE;
`,
	alter: `
`,
}
