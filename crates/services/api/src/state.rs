use lib_core::mongodb::Database;
use lib_core::repositories::location::{LocationRepository, MongoLocationRepository};
use lib_core::repositories::user::{MongoUserRepository, UserRepository};
use lib_core::services::location::LocationService;
use lib_core::services::user::UserService;

pub type ApplicationState = AppState<MongoUserRepository, MongoLocationRepository>;

#[derive(Clone)]
pub struct AppState<U: UserRepository, L: LocationRepository> {
    pub user_service: UserService<U>,
    pub location_service: LocationService<L>,
}

pub fn create_state(database: &Database) -> AppState<MongoUserRepository, MongoLocationRepository> {
    let user_repository = MongoUserRepository::new(database.collection("users"));
    let user_service = UserService::new(user_repository);

    let location_repository = MongoLocationRepository::new(database.collection("locations"));
    let location_service = LocationService::new(location_repository);

    AppState {
        user_service,
        location_service,
    }
}
