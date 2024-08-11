use axum::{
    extract::State,
    http::StatusCode,
    middleware,
    response::IntoResponse,
    routing::{get, post},
    Extension, Json, Router,
};
use lib_auth::schemas::Claims;
use lib_core::schemas::user::{CreateUserSchema, UserSchema};

use crate::{
    errors::Result,
    middlewares::auth::auth_middleware,
    state::AppState,
};

pub(crate) fn router(state: AppState) -> Router<AppState> {
    let auth_middleware = middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/register", post(register))
        .route("/me", get(profile))
        .layer(auth_middleware)
}

async fn register(
    State(state): State<AppState>,
    Extension(claims): Extension<Claims>,
    Json(user): Json<CreateUserSchema>,
) -> Result<impl IntoResponse> {
    let user = state.user_service.register(user, claims.address).await?;

    Ok((StatusCode::CREATED, Json(user)))
}

async fn profile(Extension(user): Extension<UserSchema>) -> Json<UserSchema> {
    Json(user)
}
