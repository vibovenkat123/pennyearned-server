package setup;
import (
    "fmt"
    "log"

    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
    "strings"
)

type Schema struct { create string
	drop   string
}



var defaultSchema =  Schema{

create: `
create table if not exists users (
    id text primary key not null,
    name varchar(255) not null,
    email varchar(255) not null
);

create table if not exists expenses (
    id serial primary key,
    owner_id text references users(id) not null,
    energy integer default 0,
    water integer default 0,
    personal integer default 0,
    health integer default 0,
    business integer default 0,
    rent integer default 0,
    loans integer default 0
);`,

drop: `
  drop table if exists users cascade;
  drop table if exists expenses cascade;
`,
}



func ExecMultiple(e *sqlx.DB, query string) {
    statements := strings.Split(query, "\n")
    if len(strings.Trim(statements[len(statements)-1], " \n\t\r")) == 0 {
		statements = statements[:len(statements)-1]
	}
    for _, s := range statements {
        _, err := e.Exec(s)
        if err != nil {
            log.Fatalln(err)
        }
    }
}
func Connect() (*sqlx.DB, error) {
    const (
      host     = "localhost"
      port     = 5432
      user     = "vibo"
      password = "YfWLumYuoQfCwZFLyJz8Hc27DxNfsWVfz"
      dbname   = "pennyearned"
    )

    postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

    fmt.Println("Attempting to connect...")
    db, err := sqlx.Open("postgres", postgresqlDbInfo)
    if err != nil {
        log.Fatalln(err)
    }
    err = db.Ping() 
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println("Connected!!")
    return db, err
}
func ResetToSchema(db *sqlx.DB) {
    fmt.Println("Resetting...")
    ExecMultiple(db, defaultSchema.drop)
    db.MustExec(defaultSchema.create) 
    fmt.Println("Resetted!!")
}
    
