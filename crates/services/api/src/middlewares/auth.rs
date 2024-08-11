use axum::{
    body::Body,
    extract::{Request, State},
    http::Response,
    middleware::Next,
};
use lib_auth::{errors::AuthError, jwt::decode_token};
use lib_core::errors::CoreError;

use crate::{errors::ApiError, state::AppState};

pub async fn auth_middleware(
    State(state): State<AppState>,
    mut request: Request,
    next: Next,
) -> Result<Response<Body>, ApiError> {
    let auth_header = match request.headers_mut().get(axum::http::header::AUTHORIZATION) {
        None => return Err(CoreError::from(AuthError::InvalidToken).into()),
        Some(header) => header
            .to_str()
            .map_err(|_| CoreError::from(AuthError::InvalidToken))?,
    };

    let mut header = auth_header.split_whitespace();
    let (_token_type, token) = (
        header.next(),
        header
            .next()
            .ok_or::<ApiError>(CoreError::from(AuthError::InvalidToken).into())?,
    );

    let token_data = decode_token(token).map_err(|_| CoreError::from(AuthError::InvalidToken))?;
    request.extensions_mut().insert(token_data.claims.clone());

    let user = state
        .user_service
        .find_one_user_by_address(token_data.claims.address)
        .await;
    if let Ok(user) = user {
        request.extensions_mut().insert(user);
    }

    Ok(next.run(request).await)
}
