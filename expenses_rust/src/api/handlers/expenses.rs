use std::sync::Arc;

use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use chrono::serde::ts_seconds::serialize as to_ts;
use chrono::DateTime;
use serde::Serialize;

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
        return Err(StatusCode::BAD_REQUEST)
    }
    let expense: Expense = match sqlx::query_as::<_, Expense>("SELECT * FROM expenses WHERE id=$1")
        .bind(id.clone())
        .fetch_one(&state.pool)
        .await
    {
        Ok(val) => val,
        Err(e) => match e {
            sqlx::Error::RowNotFound => {
                return Err(StatusCode::NOT_FOUND);
            }
            _ => {
                panic!("{e}");
            }
        },
    };
    Ok(Json(expense))
}
