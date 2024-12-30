use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize, Default)]
pub struct Spaceship {
    pub _id: ObjectId,
    pub name: String,
    pub user_id: ObjectId,
    pub location_id: ObjectId,
    pub flown_out_at: Option<i64>, // timestamp
    pub flying: bool,
    pub system_id: ObjectId,
    pub planet_id: Option<ObjectId>,
}
