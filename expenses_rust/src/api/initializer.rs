use std::{net::SocketAddr, sync::Arc};

use crate::api::handlers;
use axum::{
    routing::{delete, get, post, put},
    Router,
};
use sqlx::{Pool, Postgres};
use tracing::{debug, info, warn};

pub struct State {
    pub pool: Pool<Postgres>,
}
const PORT: &str = "8009";
pub async fn init(pool: Pool<Postgres>) -> Result<(), ()> {
    let state = State { pool: pool.clone() };
    let shared_state = Arc::new(state);
    if PORT.is_empty() {
        warn!(PORT, "port is empty")
    }
    info!(PORT, "Starting web server");
    axum::Server::bind(&format!("0.0.0.0:{}", PORT).parse().unwrap())
        .serve(app(shared_state).into_make_service_with_connect_info::<SocketAddr>())
        .await
        .unwrap();
    Ok(())
}

fn app(state: Arc<State>) -> Router {
    debug!("Creating the router");
    let router = Router::new()
        .route("/v1/api/expense/:id", get(handlers::expenses::get))
        .route("/v1/api/expense/:id", put(handlers::expenses::update))
        .route("/v1/api/expense/:id", delete(handlers::expenses::delete))
        .route(
            "/v1/api/expenses/user/:id",
            get(handlers::expenses::get_by_user_id),
        )
        .route("/v1/api/expense", post(handlers::expenses::new))
        .with_state(state);
    debug!("Created router; router={:?}", router);
    router
}
