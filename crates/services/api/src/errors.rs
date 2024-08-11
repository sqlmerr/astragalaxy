use axum::{
    http::{header, StatusCode},
    response::IntoResponse,
    Json,
};
use lib_core::errors::CoreError;
use lib_ton::error::Error as TonError;

pub type Result<T> = std::result::Result<T, ApiError>;

#[derive(Clone, Debug, thiserror::Error)]
pub enum ApiError {
    #[error(transparent)]
    CoreError(#[from] CoreError),
    #[error("ton error: {0}")]
    TonError(TonError),
}

impl IntoResponse for ApiError {
    fn into_response(self) -> axum::response::Response {
        let _msg = self.to_string();

        let (status_code, message) = match self {
            ApiError::CoreError(core_error) => {
                let message = core_error.to_string();
                match core_error {
                    CoreError::ServerError => (StatusCode::INTERNAL_SERVER_ERROR, message),
                    CoreError::UsernameAlreadyOccupied => (StatusCode::FORBIDDEN, message),
                    _ => (StatusCode::INTERNAL_SERVER_ERROR, message),
                }
            }
            ApiError::TonError(ton_error) => {
                let message = ton_error.0;
                (StatusCode::INTERNAL_SERVER_ERROR, message)
            }
        };

        (
            status_code,
            [(header::CONTENT_TYPE, "application/json")],
            Json(serde_json::json!({ "status_code": status_code.as_u16(), "message": message })),
        )
            .into_response()
    }
}
