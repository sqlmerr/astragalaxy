use lib_core::models::Planet;
use lib_core::mongodb::Database;
use lib_core::repositories::location::{LocationRepository, MongoLocationRepository};
use lib_core::repositories::planet::{MongoPlanetRepository, PlanetRepository};
use lib_core::repositories::spaceship::{MongoSpaceshipRepository, SpaceshipRepository};
use lib_core::repositories::system::{MongoSystemRepository, SystemRepository};
use lib_core::repositories::user::{MongoUserRepository, UserRepository};
use lib_core::services::location::LocationService;
use lib_core::services::planet::PlanetService;
use lib_core::services::spaceship::SpaceshipService;
use lib_core::services::system::SystemService;
use lib_core::services::user::UserService;

use crate::config::Config;

pub type ApplicationState = AppState<
    MongoUserRepository,
    MongoLocationRepository,
    MongoPlanetRepository,
    MongoSystemRepository,
    MongoSpaceshipRepository,
>;

#[derive(Clone)]
pub struct AppState<U, L, P, Sy, Sp>
where
    U: UserRepository,
    L: LocationRepository,
    P: PlanetRepository,
    Sy: SystemRepository,
    Sp: SpaceshipRepository,
{
    pub user_service: UserService<U>,
    pub location_service: LocationService<L>,
    pub planet_service: PlanetService<P>,
    pub system_service: SystemService<Sy>,
    pub spaceship_service: SpaceshipService<Sp>,
    pub config: Config,
}

pub fn create_state(database: &Database, config: Config) -> ApplicationState {
    let user_repository = MongoUserRepository::new(database.collection("users"));
    let user_service = UserService::new(user_repository);

    let location_repository = MongoLocationRepository::new(database.collection("locations"));
    let location_service = LocationService::new(location_repository);

    let planet_repository = MongoPlanetRepository::new(database.collection("planets"));
    let planet_service = PlanetService::new(planet_repository);

    let system_repository = MongoSystemRepository::new(database.collection("systems"));
    let system_service = SystemService::new(system_repository);

    let spaceship_repository = MongoSpaceshipRepository::new(database.collection("spaceships"));
    let spaceship_service = SpaceshipService::new(spaceship_repository);

    AppState {
        user_service,
        location_service,
        planet_service,
        system_service,
        spaceship_service,
        config,
    }
}
