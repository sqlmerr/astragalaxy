use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

use crate::models::{planet::PlanetThreat, Planet};

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct PlanetSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub system_id: ObjectId,
    pub threat: PlanetThreat,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreatePlanetSchema {
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub system_id: ObjectId,
    pub threat: PlanetThreat,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UpdatePlanetSchema {
    pub system_id: Option<ObjectId>,
    pub threat: Option<PlanetThreat>,
}

impl From<Planet> for PlanetSchema {
    fn from(value: Planet) -> Self {
        Self {
            _id: value._id,
            system_id: value.system_id,
            threat: value.threat,
        }
    }
}
