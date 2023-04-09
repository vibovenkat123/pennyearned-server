use std::env;

fn main() {
    let db_url = env::var("DATABASE_URL").unwrap_or("blah".to_string());
    println!("cargo:rustc-env=DB_URL={}", db_url);
}
