use std::str::FromStr;

use axum::{
    extract::{Path, State},
    routing::{get, post},
    Extension, Json, Router,
};
use lib_core::{
    errors::CoreError,
    mongodb::bson::oid::ObjectId,
    schemas::{spaceship::SpaceshipSchema, user::UserSchema},
};

use crate::{
    errors::Result, middlewares::auth::auth_middleware, schemas::responses::OkResponse,
    state::ApplicationState,
};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = axum::middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/:id", get(get_spaceship_by_id))
        .route("/my", get(get_my_spaceship))
        .route("/my/enter", post(enter_my_spaceship))
        .route("/my/getOut", post(get_out_of_my_spaceship))
        .layer(auth_middleware)
}

async fn get_spaceship_by_id(
    Path(id): Path<String>,
    State(state): State<ApplicationState>,
) -> Result<Json<SpaceshipSchema>> {
    let id = match ObjectId::from_str(id.as_str()) {
        Ok(id) => id,
        Err(_) => return Err(CoreError::NotFound.into()),
    };

    let spaceship = state.spaceship_service.find_one_spaceship(id).await?;

    Ok(Json(spaceship))
}

async fn get_my_spaceship(
    Extension(user): Extension<UserSchema>,
    State(state): State<ApplicationState>,
) -> Result<Json<SpaceshipSchema>> {
    let id = match user.spaceship_id {
        Some(id) => id,
        None => return Err(CoreError::PlayerHasNoSpaceship.into()),
    };

    let spaceship = state.spaceship_service.find_one_spaceship(id).await?;

    Ok(Json(spaceship))
}

async fn enter_my_spaceship(
    Extension(user): Extension<UserSchema>,
    State(state): State<ApplicationState>,
) -> Result<Json<OkResponse>> {
    let ok = state.user_service.board_spaceship(user).await.is_ok();

    Ok(Json(OkResponse::new(ok)))
}

async fn get_out_of_my_spaceship(
    Extension(user): Extension<UserSchema>,
    State(state): State<ApplicationState>,
) -> Result<Json<OkResponse>> {
    let ok = state.user_service.get_out_of_spaceship(user).await.is_ok();

    Ok(Json(OkResponse::new(ok)))
}
