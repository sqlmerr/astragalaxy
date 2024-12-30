use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

use crate::models::System;

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SystemSchema {
    #[serde(
        rename = "id",
        serialize_with = "mongodb::bson::serde_helpers::serialize_object_id_as_hex_string"
    )]
    pub _id: ObjectId,
    pub name: String,
    pub neighbours: Vec<ObjectId>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct CreateSystemSchema {
    pub name: String,
    pub neighbours: Vec<ObjectId>,
}

#[derive(Default, Debug, Clone, Deserialize, Serialize)]
pub struct UpdateSystemSchema {
    pub name: Option<String>,
    pub neighbours: Option<Vec<ObjectId>>,
}

impl From<System> for SystemSchema {
    fn from(value: System) -> Self {
        Self {
            _id: value._id,
            name: value.name,
            neighbours: value.neighbours,
        }
    }
}
