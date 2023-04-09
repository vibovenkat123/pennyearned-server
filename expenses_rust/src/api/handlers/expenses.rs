use axum::{
    extract::{self, Path, State},
    http::StatusCode,
    Json,
};
use chrono::serde::ts_seconds::serialize as to_ts;
use chrono::DateTime;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tracing::debug;
use tracing::error;
use tracing::info;
use uuid::Uuid;

#[derive(Serialize, sqlx::FromRow, Debug)]
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
#[derive(Deserialize)]
pub struct UpdateExpense {
    name: Option<String>,
    spent: Option<i32>,
}

use crate::api::initializer::State as state_struct;

fn validate_id(id: String) -> bool {
    id.len() == 36
}

#[axum_macros::debug_handler]
pub async fn get(
    State(state): State<Arc<state_struct>>,
    Path(id): Path<String>,
) -> Result<Json<Expense>, StatusCode> {
    if !validate_id(id.clone()) {
        return Err(StatusCode::BAD_REQUEST);
    }
    info!(id, "getting specific expense");
    let expense: Expense = match sqlx::query_as!(
        Expense,
        r#"
        SELECT *
          FROM expenses
          WHERE id=$1
        "#,
        id
    )
    .fetch_one(&state.pool)
    .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                error!("expense not found");
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                error!("{e}");
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };
    Ok(Json(expense))
}

#[axum_macros::debug_handler]
pub async fn delete(
    State(state): State<Arc<state_struct>>,
    Path(id): Path<String>,
) -> Result<StatusCode, StatusCode> {
    if !validate_id(id.clone()) {
        debug!("ID is not valid");
        return Err(StatusCode::BAD_REQUEST);
    }
    info!(id, "checking if expense exists to delete");
    let _ = match sqlx::query!(
        r#"
        SELECT *
          FROM expenses
          WHERE id=$1
        "#,
        id
    )
    .fetch_one(&state.pool)
    .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                error!(id, "expense not found");
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                error!("{e}");
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };
    info!(id, "deleting specific expense");
    sqlx::query!(
        r#"
        DELETE
          FROM expenses
          WHERE id=$1
        "#,
        id
    )
    .execute(&state.pool)
    .await
    .unwrap();
    Ok(StatusCode::NO_CONTENT)
}
#[axum_macros::debug_handler]
pub async fn get_by_user_id(
    State(state): State<Arc<state_struct>>,
    Path(id): Path<String>,
) -> Result<Json<Vec<Expense>>, StatusCode> {
    if !validate_id(id.clone()) {
        debug!("ID is not valid");
        return Err(StatusCode::BAD_REQUEST);
    }
    info!(id, "getting all expenses for user");
    let expense: Vec<Expense> = match sqlx::query_as!(
        Expense,
        r#"
        SELECT *
          FROM expenses
          WHERE owner_id=$1
        "#,
        id
    )
    .fetch_all(&state.pool)
    .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                error!("Failed to get expenses");
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                error!("{e}");
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };
    Ok(Json(expense))
}

#[axum_macros::debug_handler]
pub async fn new(
    State(state): State<Arc<state_struct>>,
    extract::Json(payload): extract::Json<NewExpense>,
) -> Result<StatusCode, StatusCode> {
    info!("Validating all the json payloads");
    let owner_id = match payload.owner_id {
        Some(ref val) => val,
        None => {
            error!("JSON payload not in right format");
            return Err(StatusCode::BAD_REQUEST);
        }
    };

    let name = match payload.name {
        Some(ref val) => val,
        None => {
            error!("JSON payload not in right format");
            return Err(StatusCode::BAD_REQUEST);
        }
    };

    let spent = match payload.spent {
        Some(ref val) => val,
        None => {
            error!("JSON payload not in right format");
            return Err(StatusCode::BAD_REQUEST);
        }
    };
    if !validate_id(owner_id.clone()) {
        error!("Owner ID is not in valid format");
        return Err(StatusCode::BAD_REQUEST);
    }
    info!("Generating new uuid");
    let id = Uuid::new_v4().to_string();
    debug!(id, "Generated new uuid");
    info!("Creating new expense");
    sqlx::query!(
        r#"
        INSERT INTO expenses (owner_id, name, spent, id)
          VALUES ($1, $2, $3, $4);
        "#,
        owner_id,
        name,
        spent,
        id
    )
    .execute(&state.pool)
    .await
    .unwrap();
    Ok(StatusCode::CREATED)
}

#[axum_macros::debug_handler]
pub async fn update(
    State(state): State<Arc<state_struct>>,
    Path(id): Path<String>,
    extract::Json(payload): extract::Json<UpdateExpense>,
) -> Result<StatusCode, StatusCode> {
    info!("Validating id");
    if !validate_id(id.clone()) {
        error!("ID is not in valid format");
        return Err(StatusCode::BAD_REQUEST);
    }
    info!("Getting original expense details");
    let original: Expense = match sqlx::query_as!(Expense, "SELECT * FROM expenses WHERE id=$1", id)
        .fetch_one(&state.pool)
        .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                error!("expense not found");
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                error!("{e}");
                return Err(StatusCode::INTERNAL_SERVER_ERROR);
            }
        },
    };
    info!("Checking if the json values are there or not, and updating them if so");
    let name = payload.name.clone().unwrap_or(original.name);
    let spent = payload.spent.clone().unwrap_or(original.spent);
    info!("Updating expense");
    debug!(id, "With id");
    sqlx::query!(
        r#"
        UPDATE expenses
          SET date_updated = now(), name = $1, spent = $2
          WHERE id=$3;
        "#,
        name,
        spent,
        id
    )
    .execute(&state.pool)
    .await
    .unwrap();
    Ok(StatusCode::OK)
}
