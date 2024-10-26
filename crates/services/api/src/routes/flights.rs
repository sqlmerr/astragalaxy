use axum::Router;

use crate::{middlewares::auth::auth_middleware, state::ApplicationState};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = axum::middleware::from_fn_with_state(state, auth_middleware);

    Router::new().layer(auth_middleware)
}
