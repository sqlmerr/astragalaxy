use chrono::{DateTime, Utc};
use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};
use serde_with::{serde_as, DisplayFromStr};

use crate::models::Spaceship;

#[serde_as]
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SpaceshipSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    pub name: String,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub user_id: ObjectId,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub location_id: ObjectId,
    pub flown_out_at: Option<DateTime<Utc>>,
    pub flying: bool,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub system_id: ObjectId,
    #[serde_as(as = "Option<DisplayFromStr>")]
    pub planet_id: Option<ObjectId>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreateSpaceshipSchema {
    pub name: String,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub user_id: ObjectId,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub location_id: ObjectId,
    #[serde(serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string")]
    pub system_id: ObjectId,
}

#[derive(Default, Debug, Clone, Deserialize, Serialize)]
pub struct UpdateSpaceshipSchema {
    pub name: Option<String>,
    pub user_id: Option<ObjectId>,
    pub location_id: Option<ObjectId>,
    pub flown_out_at: Option<Option<DateTime<Utc>>>,
    pub flying: Option<bool>,
    pub system_id: Option<ObjectId>,
    pub planet_id: Option<Option<ObjectId>>,
}

impl From<Spaceship> for SpaceshipSchema {
    fn from(value: Spaceship) -> Self {
        let flown_out_at = match value.flown_out_at {
            None => None,
            Some(t) => DateTime::from_timestamp(t, 0),
        };
        Self {
            _id: value._id,
            name: value.name,
            user_id: value.user_id,
            location_id: value.location_id,
            flown_out_at,
            flying: value.flying,
            system_id: value.system_id,
            planet_id: value.planet_id,
        }
    }
}
