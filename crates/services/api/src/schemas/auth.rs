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
