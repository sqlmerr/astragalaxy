use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

use crate::models::User;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UserSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    pub username: String,
    pub spaceship_id: Option<ObjectId>,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub location_id: ObjectId,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreateUserSchema {
    pub username: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UpdateUserSchema {
    pub username: Option<String>,
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
