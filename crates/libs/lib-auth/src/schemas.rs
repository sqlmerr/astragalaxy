use chrono::{Duration, Utc};
use jsonwebtoken::{DecodingKey, EncodingKey};
use lib_ton::PAYLOAD_TTL;

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
    /// Ton address
    pub address: String,
    /// Expiration
    pub exp: usize,
}

impl Claims {
    pub fn new(address: String) -> Self {
        Self {
            address,
            exp: (Utc::now() + Duration::seconds(PAYLOAD_TTL as i64)).timestamp() as usize,
        }
    }
}
