use routes::app;
use tokio::net::TcpListener;

mod routes;
pub mod state;

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok().unwrap();

    let app = app().await;
    let listener = TcpListener::bind("0.0.0.0:8000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
