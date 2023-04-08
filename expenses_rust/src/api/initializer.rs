use std::sync::Arc;

use crate::api::handlers;
use axum::{routing::get, Router};
use sqlx::{Pool, Postgres};

pub struct State {
    pub pool: Pool<Postgres>,
}

pub async fn init(pool: Pool<Postgres>) -> Result<(), ()> {
    let state = State { pool: pool.clone() };
    let shared_state = Arc::new(state);
    axum::Server::bind(&"0.0.0.0:8009".parse().unwrap())
        .serve(app(shared_state).into_make_service())
        .await
        .unwrap();
    Ok(())
}

fn app(state: Arc<State>) -> Router {
    Router::new()
        .route("/v1/api/expense/:id", get(handlers::expenses::get))
        .route(
            "/v1/api/expenses/user/:id",
            get(handlers::expenses::get_by_user_id),
        )
        .with_state(state)
}
