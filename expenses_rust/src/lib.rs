use std::process;

mod api;
mod db;
pub async fn initialize() -> Result<(), ()> {
    let pool = match db::setup::setup().await {
        Ok(val) => val,
        Err(e) => {
            eprintln!("{e}");
            process::exit(1)
        }
    };
    api::initializer::init(pool).await?;
    Ok(())
}
