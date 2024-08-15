use axum::{
    extract::State,
    http::StatusCode,
    middleware,
    routing::{get, post},
    Extension, Json, Router,
};
use lib_auth::schemas::Claims;
use lib_core::schemas::user::{CreateUserSchema, UserSchema};

use crate::{errors::Result, middlewares::auth::auth_middleware, state::ApplicationState};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/register", post(register))
        .route("/me", get(profile))
        .layer(auth_middleware)
}

async fn register(
    State(state): State<ApplicationState>,
    Extension(claims): Extension<Claims>,
    Json(user): Json<CreateUserSchema>,
) -> Result<(StatusCode, Json<UserSchema>)> {
    let location = state
        .location_service
        .find_one_location_by_code("space_station".to_string())
        .await?;
    let user = state
        .user_service
        .register(user, claims.address, location._id)
        .await?;

    Ok((StatusCode::CREATED, Json(user)))
}

async fn profile(Extension(user): Extension<UserSchema>) -> Json<UserSchema> {
    Json(user)
}
