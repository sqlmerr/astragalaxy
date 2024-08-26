use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize, Default)]
pub struct User {
    pub _id: ObjectId,
    pub username: String,
    pub discord_id: Option<i64>,
    pub password: Option<String>,
    pub spaceship_id: Option<ObjectId>,
    pub location_id: ObjectId,
    pub x: i64,
    pub y: i64,
}
