use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};
use serde_with::serde_as;

use serde_with::DisplayFromStr;

use crate::models::User;

#[serde_as]
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UserSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    pub username: String,
    pub telegram_id: i64,
    #[serde_as(as = "Option<DisplayFromStr>")]
    pub spaceship_id: Option<ObjectId>,
    pub in_spaceship: bool,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub location_id: ObjectId,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub system_id: ObjectId,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreateUserSchema {
    pub username: String,
    pub telegram_id: i64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UpdateUserSchema {
    pub username: Option<String>,
    pub spaceship_id: Option<ObjectId>,
}

impl From<User> for UserSchema {
    fn from(value: User) -> Self {
        Self {
            _id: value._id,
            username: value.username,
            spaceship_id: value.spaceship_id,
            telegram_id: value.telegram_id,
            in_spaceship: value.in_spaceship,
            location_id: value.location_id,
            system_id: value.system_id,
        }
    }
}
