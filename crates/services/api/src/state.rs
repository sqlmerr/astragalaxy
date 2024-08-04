use lib_core::services::user::UserService;

#[derive(Clone)]
pub struct AppState {
    pub user_service: UserService,
}
