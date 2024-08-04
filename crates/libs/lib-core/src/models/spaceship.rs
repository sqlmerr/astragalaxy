use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Spaceship {
    pub _id: ObjectId,
    pub name: String,
    pub user_id: ObjectId,
}
