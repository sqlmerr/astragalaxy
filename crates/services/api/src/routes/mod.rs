use axum::Router;
use lib_core::{
    create_mongodb_client, repositories::user::UserRepository, services::user::UserService,
};

use crate::state::AppState;

pub async fn app() -> Router {
    let client =
        create_mongodb_client(std::env::var("MONGODB_URI").expect("MONGODB_URI must be set")).await;
    let database = client.database("astragalaxy");

    let user_repository = UserRepository::new(database.collection("users"));
    let user_service = UserService::new(user_repository);

    let state = AppState { user_service };

    Router::new().with_state(state)
}
