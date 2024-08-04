pub mod errors;
pub mod models;
pub mod repositories;
pub mod schemas;
pub mod services;

pub(crate) use errors::Result;
use models::User;
pub use mongodb;
use mongodb::{bson::doc, options::IndexOptions, Client, IndexModel};

pub async fn create_mongodb_client(uri: String) -> Client {
    Client::with_uri_str(uri.as_str()).await.unwrap()
}

pub async fn startup(client: &Client) {
    let options = IndexOptions::builder().unique(true).build();
    let model = IndexModel::builder()
        .keys(doc! {"username": 1})
        .options(options)
        .build();
    client
        .database("astragalaxy")
        .collection::<User>("users")
        .create_index(model)
        .await
        .expect("error while creating index");
}
