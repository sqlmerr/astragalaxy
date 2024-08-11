use lib_core::{
    create_mongodb_client, repositories::location::LocationRepository,
    schemas::location::CreateLocationSchema, services::location::LocationService, startup,
};

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok().unwrap();

    let client = create_mongodb_client(std::env::var("MONGODB_URI").unwrap()).await;
    let database = client.database("astragalaxy");
    // startup(&client).await

    let location_repo = LocationRepository::new(database.collection("locations"));
    let location_service = LocationService::new(location_repo);

    let location = location_service
        .create_location(CreateLocationSchema {
            code: "space_station".to_string(),
            multiplayer: true,
        })
        .await
        .expect("Failed to create space_station location");
    println!("Created location: {:?}", location);
}
