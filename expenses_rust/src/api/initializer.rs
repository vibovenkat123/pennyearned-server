use http::{
    header::{ACCEPT, ACCEPT_ENCODING, AUTHORIZATION, CONTENT_TYPE, ORIGIN},
    Request,
};
use std::sync::Arc;
use tower_http::{compression::CompressionLayer, cors::CorsLayer, trace::TraceLayer};

use crate::api::handlers;
use axum::{
    body::Body,
    routing::{delete, get, post, put},
    Router,
};
use sqlx::{Pool, Postgres};
use tracing::{debug, info};

pub struct State {
    pub pool: Pool<Postgres>,
}

pub async fn init(pool: Pool<Postgres>) -> Result<(), ()> {
    let state = State { pool: pool.clone() };
    let shared_state = Arc::new(state);
    #[cfg(debug_assertions)]
    {
        let addr = std::net::SocketAddr::from(([127, 0, 0, 1], 8009));
        info!("addr={}, {}", addr.to_string(), "Starting web server");
        axum::Server::bind(&addr)
            .serve(app(shared_state).into_make_service())
            .await
            .unwrap();
    }

    #[cfg(not(debug_assertions))]
    {
        info!("Starting web server on lambda");
        let app = tower::ServiceBuilder::new()
            .layer(axum_aws_lambda::LambdaLayer::default())
            .service(app(shared_state));
        lambda_http::run(app).await.unwrap();
    }
    Ok(())
}

fn app(state: Arc<State>) -> Router {
    debug!("Creating the layers");
    let trace_layer =
        TraceLayer::new_for_http().on_request(|_: &Request<Body>, _: &tracing::Span| {
            tracing::info!(message = "Recieved request")
        });

    let cors_layer = CorsLayer::new()
        .allow_headers(vec![
            ACCEPT,
            ACCEPT_ENCODING,
            AUTHORIZATION,
            CONTENT_TYPE,
            ORIGIN,
        ])
        .allow_methods(tower_http::cors::Any)
        .allow_origin(tower_http::cors::Any);

    let compression_layer = CompressionLayer::new().gzip(true).deflate(true);
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
        .layer(cors_layer)
        .layer(trace_layer)
        .layer(compression_layer)
        .with_state(state);
    debug!("Created router; router={:?}", router);
    router
}
