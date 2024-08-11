use lib_auth::errors::AuthError;
use thiserror::Error;

pub(crate) type Result<T> = std::result::Result<T, CoreError>;

#[derive(Clone, Error, Debug)]
pub enum CoreError {
    #[error("This username is already occupied")]
    UsernameAlreadyOccupied,
    #[error("Can't create user")]
    CantCreateUser,
    #[error("Server error")]
    ServerError,
    #[error("Entity not found")]
    NotFound,
    #[error(transparent)]
    AuthError(#[from] AuthError),
    #[error("You already have an account with this address")]
    AddressAlreadyOccupied,
}
