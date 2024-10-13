use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

use crate::models::Spaceship;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SpaceshipSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    pub name: String,
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub user_id: ObjectId,
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub location_id: ObjectId,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreateSpaceshipSchema {
    pub name: String,
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub user_id: ObjectId,
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub location_id: ObjectId,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UpdateSpaceshipSchema {
    pub name: Option<String>,
    pub user_id: Option<ObjectId>,
    pub location_id: Option<ObjectId>,
}

impl From<Spaceship> for SpaceshipSchema {
    fn from(value: Spaceship) -> Self {
        Self {
            _id: value._id,
            name: value.name,
            user_id: value.user_id,
            location_id: value.location_id,
        }
    }
}
