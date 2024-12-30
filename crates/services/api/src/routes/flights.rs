use axum::{extract::State, routing::post, Extension, Json, Router};
use lib_core::{errors::CoreError, schemas::user::UserSchema};

use crate::{
    errors::Result,
    middlewares::auth::auth_middleware,
    schemas::{responses::OkResponse, spaceship::SpaceshipFlySchema},
    state::ApplicationState,
};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = axum::middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/planet", post(fly_to_planet))
        .layer(auth_middleware)
}

async fn fly_to_planet(
    State(state): State<ApplicationState>,
    Extension(user): Extension<UserSchema>,
    Json(data): Json<SpaceshipFlySchema>,
) -> Result<Json<OkResponse>> {
    let spaceship_id = match user.spaceship_id {
        None => return Err(CoreError::PlayerHasNoSpaceship.into()),
        Some(s) => s,
    };

    state
        .spaceship_service
        .fly(spaceship_id, data.planet_id)
        .await?;
    println!("aaa");

    Ok(Json(OkResponse::new(true)))
}
