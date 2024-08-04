use jsonwebtoken::{decode, encode, errors::Result, Header, TokenData, Validation};

use once_cell::sync::Lazy;

use crate::schemas::{Claims, Keys};

static KEYS: Lazy<Keys> = Lazy::new(|| {
    let secret = std::env::var("JWT_SECRET").expect("JWT_SECRET must be set");
    Keys::new(secret.as_bytes())
});

pub fn create_token(claims: &Claims) -> Result<String> {
    encode(&Header::default(), claims, &KEYS.encoding)
}

pub fn decode_token(token: &str) -> Result<TokenData<Claims>> {
    decode::<Claims>(token, &KEYS.decoding, &Validation::default())
}
