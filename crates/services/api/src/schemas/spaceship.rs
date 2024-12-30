use lib_core::mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
pub struct SpaceshipRenameSchema {
    pub name: String,
}

#[derive(Serialize, Deserialize)]
pub struct SpaceshipFlySchema {
    #[serde(
        serialize_with = "lib_core::mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub planet_id: ObjectId,
}
