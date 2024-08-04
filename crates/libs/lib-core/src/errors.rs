use thiserror::Error;

pub(crate) type Result<T> = std::result::Result<T, CoreError>;

#[derive(Error, Debug)]
pub enum CoreError {
    #[error("This username is already occupied")]
    UsernameAlreadyOccupied,
    #[error("Can't create user")]
    CantCreateUser,
    #[error("Server error")]
    ServerError,
    #[error("User not found")]
    UserNotFound,
}
