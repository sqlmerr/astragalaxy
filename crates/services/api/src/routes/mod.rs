mod auth;

use axum::{http::StatusCode, response::IntoResponse, Json, Router};
use lib_core::{
    create_mongodb_client, mongodb::Database, repositories::user::UserRepository,
    services::user::UserService,
};
use serde_json::json;

use crate::state::AppState;

pub async fn app(database: Database) -> Router {
    let user_repository = UserRepository::new(database.collection("users"));
    let user_service = UserService::new(user_repository);

    let state = AppState { user_service };

    Router::new()
        .nest("/auth", auth::router(state.clone()))
        .layer(tower_http::trace::TraceLayer::new_for_http())
        .fallback(handler_404)
        .with_state(state)
}

async fn handler_404() -> impl IntoResponse {
    (StatusCode::NOT_FOUND, Json(json!({"message": "Not found"})))
}
