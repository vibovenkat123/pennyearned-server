use tracing_subscriber;
use tracing::Level;
mod api;
mod db;
pub async fn initialize() -> Result<(), ()> {
    #[cfg(debug_assertions)]
    {
        tracing_subscriber::fmt()
            .with_max_level(Level::DEBUG)
            .init();
    }
    #[cfg(not(debug_assertions))]
    {
        tracing_subscriber::fmt()
            .with_ansi(false)
            .with_max_level(Level::INFO)
            .init();
    }
    let pool = match db::setup::setup().await {
        Ok(val) => val,
        Err(e) => {
            panic!("{e}");
        }
    };
    api::initializer::init(pool).await?;
    Ok(())
}
