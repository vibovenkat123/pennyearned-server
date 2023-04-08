use std::env;

fn main() {
    let db_url = env::var("DB_URL").unwrap_or("blah".to_string());
    println!("cargo:rustc-env=DB_URl={}", db_url);
}
