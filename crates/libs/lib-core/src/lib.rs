pub mod errors;
pub mod models;
pub mod repositories;
pub mod schemas;
pub mod services;

pub(crate) use errors::Result;
use mongodb::Client;

pub async fn create_mongodb_client(uri: String) -> Client {
    Client::with_uri_str(uri.as_str()).await.unwrap()
}
