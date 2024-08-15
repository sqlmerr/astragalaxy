use crate::errors::ApiError;
use crate::middlewares::auth::auth_middleware;
use crate::state::ApplicationState;
use axum::routing::{get, post};
use axum::{middleware, Extension, Json, Router};
use lib_auth::jwt::create_token;
use lib_auth::schemas::Claims;
use lib_core::errors::CoreError;
use lib_ton::{
    check_ton_proof, generate_ton_proof_payload, CheckProofPayload, CheckTonProof,
    GenerateTonProofPayload, WalletAddress,
};

pub(crate) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/generate-payload", post(generate_ton_proof))
        .route("/check-proof", post(check_proof))
        .route(
            "/get-account-info",
            get(get_account_info).layer(auth_middleware),
        )
}

async fn generate_ton_proof() -> Result<Json<GenerateTonProofPayload>, ApiError> {
    let secret = std::env::var("TON_SECRET").expect("TON_SECRET must be set");
    let payload = generate_ton_proof_payload(secret)
        .await
        .map_err(|e| ApiError::TonError(e))?;
    Ok(Json(payload))
}

async fn check_proof(Json(body): Json<CheckProofPayload>) -> Result<Json<CheckTonProof>, ApiError> {
    let secret = std::env::var("TON_SECRET").expect("TON_SECRET must be set");
    let domain = std::env::var("TON_DOMAIN").expect("TON_DOMAIN must be set");
    let address = check_ton_proof(body, secret, domain)
        .await
        .map_err(|e| ApiError::TonError(e))?;
    let token = create_token(&Claims::new(address)).map_err(|_| CoreError::ServerError)?;

    Ok(Json(CheckTonProof { token }))
}

async fn get_account_info(
    Extension(claims): Extension<Claims>,
) -> Result<Json<WalletAddress>, ApiError> {
    Ok(Json(WalletAddress {
        address: claims.address,
    }))
}
