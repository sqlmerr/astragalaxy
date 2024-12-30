use axum::{extract::State, routing::post, Json, Router};
use lib_core::schemas::planet::{CreatePlanetSchema, PlanetSchema};

use crate::{
    errors::Result,
    middlewares::auth::{auth_middleware, protection_middleware},
    state::ApplicationState,
};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    // let auth_middleware = axum::middleware::from_fn_with_state(state, auth_middleware);
    let protection_middleware = axum::middleware::from_fn_with_state(state, protection_middleware);

    Router::new().route("/", post(create_planet).layer(protection_middleware))
}

async fn create_planet(
    State(state): State<ApplicationState>,
    Json(data): Json<CreatePlanetSchema>,
) -> Result<Json<PlanetSchema>> {
    let response = state.planet_service.create_planet(data).await?;

    Ok(Json(response))
}
