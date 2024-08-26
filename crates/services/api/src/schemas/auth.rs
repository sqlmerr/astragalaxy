use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct AuthBody {
    pub access_token: String,
    pub token_type: String,
}

#[derive(Serialize, Deserialize)]
pub struct AuthPayload {
    pub username: String,
    pub password: String,
}

impl AuthBody {
    pub fn new(access_token: String) -> Self {
        Self {
            access_token,
            token_type: "Bearer".to_string(),
        }
    }
}

#[derive(Serialize, Deserialize)]
pub struct DiscordAuthPayload {
    pub code: String,
}

#[derive(Serialize, Deserialize)]
pub struct DiscordAuthBody {
    pub access_token: String,
    pub token_type: String,
    pub expires_in: i64,
    pub refresh_token: String,
    pub scope: String,
}
