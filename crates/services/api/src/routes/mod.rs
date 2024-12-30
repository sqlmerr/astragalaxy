mod auth;
mod flights;
mod planets;
mod players;
mod spaceships;
mod systems;

use axum::http::HeaderValue;
use axum::routing::get;
use axum::{http::StatusCode, response::IntoResponse, Json, Router};
use lib_core::mongodb::Database;
use serde_json::json;
use tower_http::cors::Any;

use crate::config::Config;
use crate::state::{create_state, ApplicationState};

pub async fn app(state: ApplicationState, database: Database, config: Config) -> Router {
    // let state = create_state(&database, config.clone());

    Router::new()
        .route(
            "/",
            get(|| async move {
                (
                    StatusCode::OK,
                    Json(json!({"message": "Hello World!", "ok": true})),
                )
            }),
        )
        .nest("/auth", auth::router(state.clone()))
        .nest("/players", players::router(state.clone()))
        .nest("/spaceships", spaceships::router(state.clone()))
        .nest("/flights", flights::router(state.clone()))
        .nest("/planets", planets::router(state.clone()))
        .nest("/systems", systems::router(state.clone()))
        .layer(tower_http::trace::TraceLayer::new_for_http())
        .layer(
            tower_http::cors::CorsLayer::new()
                .allow_origin(config.domain.parse::<HeaderValue>().unwrap())
                .allow_headers(Any)
                .allow_methods(Any),
        )
        .fallback(handler_404)
        .with_state(state)
}

async fn handler_404() -> impl IntoResponse {
    (StatusCode::NOT_FOUND, Json(json!({"message": "Not found"})))
}
