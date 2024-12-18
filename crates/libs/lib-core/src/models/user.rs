use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize, Default)]
pub struct User {
    pub _id: ObjectId,
    pub username: String,
    pub discord_id: Option<i64>,
    pub telegram_id: i64,
    pub spaceship_id: Option<ObjectId>,
    pub in_spaceship: bool,
    pub location_id: ObjectId,
    pub system_id: ObjectId,
    pub x: i64,
    pub y: i64,
    pub token: String,
}
