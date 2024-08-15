mod auth;
use axum::http::HeaderValue;
use axum::routing::get;
use axum::{http::StatusCode, response::IntoResponse, Json, Router};
use lib_core::mongodb::Database;
use serde_json::json;
use tower_http::cors::Any;

use crate::state::create_state;

pub async fn app(database: Database) -> Router {
    let domain = format!("https://{}", std::env::var("TON_DOMAIN").unwrap());
    let state = create_state(&database);

    Router::new()
        .route(
            "/",
            get(|| async move { (StatusCode::OK, Json(json!({"message": "Hello World!"}))) }),
        )
        .nest("/auth", auth::router(state.clone()))
        .layer(tower_http::trace::TraceLayer::new_for_http())
        .layer(
            tower_http::cors::CorsLayer::new()
                .allow_origin(domain.parse::<HeaderValue>().unwrap())
                .allow_headers(Any)
                .allow_methods(Any),
        )
        .fallback(handler_404)
        .with_state(state)
}

async fn handler_404() -> impl IntoResponse {
    (StatusCode::NOT_FOUND, Json(json!({"message": "Not found"})))
}
