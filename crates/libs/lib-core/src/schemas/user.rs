use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

use crate::models::User;

#[derive(Clone, Deserialize, Serialize)]
pub struct UserSchema {
    #[serde(rename = "id")]
    pub _id: ObjectId,
    pub username: String,
    pub spaceship_id: Option<ObjectId>,
    pub location_id: ObjectId,
}

#[derive(Clone, Deserialize, Serialize)]
pub struct CreateUserSchema {
    pub username: String,
    pub password: String,
}

#[derive(Clone, Deserialize, Serialize)]
pub struct UpdateUserSchema {
    pub username: Option<String>,
    pub password: Option<String>,
}

impl From<User> for UserSchema {
    fn from(value: User) -> Self {
        Self {
            _id: value._id,
            username: value.username,
            spaceship_id: value.spaceship_id,
            location_id: value.location_id,
        }
    }
}
