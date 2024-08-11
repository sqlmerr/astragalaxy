mod auth;
mod ton_proof;

use axum::http::HeaderValue;
use axum::{http::StatusCode, response::IntoResponse, Json, Router};
use lib_core::{
    mongodb::Database, repositories::user::UserRepository,
    services::user::UserService,
};
use serde_json::json;
use tower_http::cors::Any;

use crate::state::AppState;

pub async fn app(database: Database) -> Router {
    let user_repository = UserRepository::new(database.collection("users"));
    let user_service = UserService::new(user_repository);

    let state = AppState { user_service };

    Router::new()
        .nest("/auth", auth::router(state.clone()))
        .nest("/ton-proof", ton_proof::router(state.clone()))
        .layer(tower_http::trace::TraceLayer::new_for_http())
        .layer(
            tower_http::cors::CorsLayer::new()
                .allow_origin(
                    "https://astragalaxy.vercel.app"
                        .parse::<HeaderValue>()
                        .unwrap(),
                )
                .allow_headers(Any)
                .allow_methods(Any),
        )
        .fallback(handler_404)
        .with_state(state)
}

async fn handler_404() -> impl IntoResponse {
    (StatusCode::NOT_FOUND, Json(json!({"message": "Not found"})))
}
