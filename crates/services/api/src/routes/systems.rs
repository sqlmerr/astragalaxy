use axum::{extract::State, routing::post, Json, Router};
use lib_core::schemas::system::{CreateSystemSchema, SystemSchema};

use crate::{
    errors::Result,
    middlewares::auth::{auth_middleware, protection_middleware},
    state::ApplicationState,
};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    // let auth_middleware = axum::middleware::from_fn_with_state(state, auth_middleware);
    let protection_middleware = axum::middleware::from_fn_with_state(state, protection_middleware);

    Router::new().route("/", post(create_system).layer(protection_middleware))
}

async fn create_system(
    State(state): State<ApplicationState>,
    Json(data): Json<CreateSystemSchema>,
) -> Result<Json<SystemSchema>> {
    let system = state.system_service.create_system(data).await?;
    Ok(Json(system))
}
