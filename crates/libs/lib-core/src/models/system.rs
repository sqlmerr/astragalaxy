use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Default, Clone, Debug, Deserialize, Serialize)]
pub struct System {
    pub _id: ObjectId,
    pub name: String,
    pub neighbours: Vec<ObjectId>,
}
