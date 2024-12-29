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
    #[error("Player hasn't a spaceship")]
    PlayerHasNoSpaceship,
    #[error("Player is already in a spaceship")]
    PlayerAlreadyInSpaceship,
    #[error("Spaceship is already flying")]
    SpaceshipAlreadyFlying,
    #[error("Planet is in another system")]
    PlanetIsInAnotherSystem,
    #[error("Spaceship is already in this planet")]
    SpaceshipIsAlreadyInThisPlanet,
}
