use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

use crate::models::location::Location;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct LocationSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    pub code: String,
    pub multiplayer: bool,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreateLocationSchema {
    pub code: String,
    pub multiplayer: bool,
}

impl From<Location> for LocationSchema {
    fn from(value: Location) -> Self {
        Self {
            _id: value._id,
            code: value.code,
            multiplayer: value.multiplayer,
        }
    }
}
