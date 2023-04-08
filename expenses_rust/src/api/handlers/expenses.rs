use axum::{
    extract::{self, ConnectInfo, Path, State},
    http::StatusCode,
    Json,
};
use chrono::serde::ts_seconds::serialize as to_ts;
use chrono::DateTime;
use serde::{Deserialize, Serialize};
use std::{net::SocketAddr, sync::Arc};
use tracing::debug;
use tracing::error;
use tracing::info;
use uuid::Uuid;

#[derive(Serialize, sqlx::FromRow)]
pub struct Expense {
    id: String,
    owner_id: String,
    name: String,
    #[serde(serialize_with = "to_ts")]
    date_created: DateTime<chrono::Utc>,
    #[serde(serialize_with = "to_ts")]
    date_updated: DateTime<chrono::Utc>,
    spent: i32,
}

#[derive(Deserialize)]
pub struct NewExpense {
    owner_id: Option<String>,
    name: Option<String>,
    spent: Option<i32>,
}

use crate::api::initializer::State as state_struct;

fn validate_id(id: String) -> bool {
    id.len() == 36
}

#[axum_macros::debug_handler]
pub async fn get(
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    State(state): State<Arc<state_struct>>,
    Path(id): Path<String>,
) -> Result<Json<Expense>, StatusCode> {
    info!("Incoming request from address {}", addr);
    if !validate_id(id.clone()) {
        return Err(StatusCode::BAD_REQUEST);
    }
    info!(id, "addr={}, {}", addr, "getting specific expense");
    let expense: Expense = match sqlx::query_as::<_, Expense>("SELECT * FROM expenses WHERE id=$1")
        .bind(id.clone())
        .fetch_one(&state.pool)
        .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                error!("addr={}, {}", addr, "expense not found");
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                error!("addr={}, {e}", addr);
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };
    Ok(Json(expense))
}

#[axum_macros::debug_handler]
pub async fn get_by_user_id(
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    State(state): State<Arc<state_struct>>,
    Path(id): Path<String>,
) -> Result<Json<Vec<Expense>>, StatusCode> {
    info!("Incoming request from address {}", addr);
    if !validate_id(id.clone()) {
        return Err(StatusCode::BAD_REQUEST);
    }
    info!(id, "addr={}, {}", addr, "getting all expenses for user");
    let expense: Vec<Expense> = match sqlx::query_as("SELECT * FROM expenses where owner_id=$1")
        .bind(id.clone())
        .fetch_all(&state.pool)
        .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                error!("addr={}, {}", addr, "Failed to get expenses");
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                error!("addr={}, {e}", addr);
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };
    Ok(Json(expense))
}

#[axum_macros::debug_handler]
pub async fn new(
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    State(state): State<Arc<state_struct>>,
    extract::Json(payload): extract::Json<NewExpense>,
) -> Result<StatusCode, StatusCode> {
    info!("Incoming request from address {}", addr);
    debug!("Validating all the json payloads");
    let owner_id = match payload.owner_id {
        Some(ref val) => val,
        None => {
            error!("addr={}, {}", addr, "JSON payload not in right format");
            return Err(StatusCode::BAD_REQUEST);
        }
    };

    let name = match payload.name {
        Some(ref val) => val,
        None => {
            error!("addr={}, {}", addr, "JSON payload not in right format");
            return Err(StatusCode::BAD_REQUEST);
        }
    };

    let spent = match payload.spent {
        Some(ref val) => val,
        None => {
            error!("addr={}, {}", addr, "JSON payload not in right format");
            return Err(StatusCode::BAD_REQUEST);
        }
    };
    if !validate_id(owner_id.clone()) {
        return Err(StatusCode::BAD_REQUEST);
    }
    let id = Uuid::new_v4().to_string();
    debug!(id, "Generated new uuid");
    info!("addr={}, {}", addr, "Creating new expense");
    sqlx::query("INSERT INTO expenses (owner_id, name, spent, id) VALUES ($1, $2, $3, $4)")
        .bind(owner_id)
        .bind(name)
        .bind(spent)
        .bind(id)
        .execute(&state.pool)
        .await
        .unwrap();
    Ok(StatusCode::CREATED)
}
