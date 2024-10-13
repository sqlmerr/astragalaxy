use std::str::FromStr;

use axum::{
    extract::{Path, State},
    routing::get,
    Json, Router,
};
use lib_core::{errors::CoreError, mongodb::bson::oid::ObjectId, schemas::user::UserSchema};

use crate::{errors::Result, middlewares::auth::auth_middleware, state::ApplicationState};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = axum::middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/:id", get(get_player_by_id))
        .layer(auth_middleware)
}

async fn get_player_by_id(
    Path(id): Path<String>,
    State(state): State<ApplicationState>,
) -> Result<Json<UserSchema>> {
    let id = match ObjectId::from_str(id.as_str()) {
        Ok(id) => id,
        Err(_) => return Err(CoreError::ServerError.into()),
    };

    let user = state.user_service.find_one_user(id).await?;

    Ok(Json(user))
}
