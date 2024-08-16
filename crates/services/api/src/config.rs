use serde::Deserialize;

#[derive(Deserialize, Clone, Debug, Default)]
pub struct Config {
    pub mongodb_uri: String,
    pub jwt_secret: String,
    pub discord_client_id: String,
    pub discord_client_secret: String,
    pub domain: String,
}

impl Config {
    pub fn from_env() -> Config {
        dotenvy::dotenv().ok();

        config::Config::builder()
            .add_source(config::Environment::default())
            .build()
            .unwrap()
            .try_deserialize()
            .unwrap()
    }
}
