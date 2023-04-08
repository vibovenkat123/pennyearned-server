#[tokio::main]
async fn main() -> Result<(), ()> {
    expenses::initialize().await?;
    Ok(())
}
