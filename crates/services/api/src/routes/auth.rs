use axum::{
    extract::State,
    http::StatusCode,
    middleware,
    response::IntoResponse,
    routing::{get, post},
    Extension, Json, Router,
};
use lib_core::schemas::user::{CreateUserSchema, UserSchema};

use crate::{
    errors::Result,
    middlewares::auth::auth_middleware,
    schemas::auth::{AuthBody, AuthPayload},
    state::AppState,
};

pub(crate) fn router(state: AppState) -> Router<AppState> {
    let auth_middleware = middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .route("/me", get(profile).layer(auth_middleware))
}

async fn register(
    State(state): State<AppState>,
    Json(user): Json<CreateUserSchema>,
) -> Result<impl IntoResponse> {
    let user = state.user_service.register(user).await?;

    Ok((StatusCode::CREATED, Json(user)))
}

async fn login(
    State(state): State<AppState>,
    Json(payload): Json<AuthPayload>,
) -> Result<Json<AuthBody>> {
    let token = state
        .user_service
        .login(payload.username, payload.password)
        .await?;

    Ok(Json(AuthBody::new(token)))
}

async fn profile(Extension(user): Extension<UserSchema>) -> Json<UserSchema> {
    Json(user)
}
