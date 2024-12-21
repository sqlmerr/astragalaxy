use lib_core::{create_mongodb_client, startup};
use routes::app;
use tokio::net::TcpListener;

pub(crate) mod config;
pub mod errors;
pub mod middlewares;
mod routes;
pub(crate) mod schemas;
pub(crate) mod state;

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok().unwrap();

    let filter = tracing_subscriber::filter::EnvFilter::default()
        .add_directive(tracing::Level::INFO.into())
        .add_directive("tower_http=trace".parse().unwrap());

    let subscriber = tracing_subscriber::fmt()
        .compact()
        .with_file(true)
        .with_line_number(true)
        .with_env_filter(filter)
        .finish();
    tracing::subscriber::set_global_default(subscriber).unwrap();

    let config = config::Config::from_env();

    let client = create_mongodb_client(config.clone().mongodb_uri).await;
    startup(&client).await;
    let database = client.database("astragalaxy");

    let app = app(database, config).await;
    let listener = TcpListener::bind("0.0.0.0:8000").await.unwrap();
    tracing::info!("Starting api on http://0.0.0.0:8000");
    axum::serve(listener, app).await.unwrap();
}
