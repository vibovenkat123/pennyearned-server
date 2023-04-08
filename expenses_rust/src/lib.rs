use tracing_subscriber;
mod api;
mod db;
pub async fn initialize() -> Result<(), ()> {
    tracing_subscriber::fmt::init();
    let pool = match db::setup::setup().await {
        Ok(val) => val,
        Err(e) => {
            panic!("{e}");
        }
    };
    api::initializer::init(pool).await?;
    Ok(())
}
