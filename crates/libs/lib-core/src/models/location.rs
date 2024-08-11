use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Default, Clone, Debug, Deserialize, Serialize)]
pub struct Location {
    pub _id: ObjectId,
    pub code: String,
    pub multiplayer: bool,
}
