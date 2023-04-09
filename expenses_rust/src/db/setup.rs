use sqlx::{postgres, Pool, Postgres};
pub async fn setup() -> Result<Pool<Postgres>, String> {
    let db_url: &'static str = option_env!("DB_URL").unwrap();
    let pool = match postgres::PgPoolOptions::new()
        .max_connections(5)
        .connect(db_url)
        .await
    {
        Ok(val) => val,
        Err(e) => {
            return Err(format!("Failed to connect to postgres db: {e}"));
        }
    };
    Ok(pool)
}
