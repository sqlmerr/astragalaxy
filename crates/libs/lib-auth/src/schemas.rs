use chrono::{Duration, Utc};
use jsonwebtoken::{DecodingKey, EncodingKey};

pub struct Keys {
    pub encoding: EncodingKey,
    pub decoding: DecodingKey,
}

impl Keys {
    pub fn new(secret: &[u8]) -> Self {
        Self {
            encoding: EncodingKey::from_secret(secret),
            decoding: DecodingKey::from_secret(secret),
        }
    }
}

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Claims {
    /// Username
    pub sub: String,
    /// Expiration
    pub exp: usize,
}

impl Claims {
    pub fn new(username: String) -> Self {
        Self {
            sub: username,
            exp: (Utc::now() + Duration::days(1)).timestamp() as usize,
        }
    }
}
