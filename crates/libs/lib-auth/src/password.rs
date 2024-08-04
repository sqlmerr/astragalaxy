use bcrypt::{hash, verify, DEFAULT_COST};

pub fn hash_password(password: String) -> String {
    hash(password, DEFAULT_COST).unwrap().to_string()
}

pub fn verify_password(password: String, hashed_password: String) -> bool {
    verify(password, hashed_password.as_str()).unwrap()
}
