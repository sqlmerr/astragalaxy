#[derive(Clone, Debug, thiserror::Error)]
pub enum AuthError {
    #[error("Wrong credentials")]
    WrongCredentials,
    #[error("Missing credentials")]
    MissingCredentials,
    #[error("Token creation error")]
    TokenCreation,
    #[error("Invalid token")]
    InvalidToken,
}

pub(crate) type Result<T> = std::result::Result<T, AuthError>;
