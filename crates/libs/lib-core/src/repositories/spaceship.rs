use async_trait::async_trait;
use chrono::{DateTime, Utc};
use futures::TryStreamExt;
use mongodb::{
    bson::{doc, oid::ObjectId, Bson, Document},
    Collection,
};
use serde::Serialize;

use crate::{errors::CoreError, models::Spaceship, Result};

#[derive(Clone)]
pub struct MongoSpaceshipRepository {
    collection: Collection<Spaceship>,
}

pub struct CreateSpaceshipDTO {
    pub name: String,
    pub user_id: ObjectId,
    pub location_id: ObjectId,
    pub system_id: ObjectId,
}

#[derive(Default, Serialize)]
pub struct UpdateSpaceshipDTO {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub user_id: Option<ObjectId>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub location_id: Option<ObjectId>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub flown_out_at: Option<Option<DateTime<Utc>>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub flying: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub system_id: Option<ObjectId>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub planet_id: Option<Option<ObjectId>>,
}

#[async_trait]
pub trait SpaceshipRepository {
    async fn create(&self, data: CreateSpaceshipDTO) -> Result<ObjectId>;
    async fn find_one(&self, oid: ObjectId) -> Result<Option<Spaceship>>;
    async fn find_all(&self, filter: Document) -> Result<Vec<Spaceship>>;
    async fn delete(&self, oid: ObjectId) -> Result<()>;
    async fn update(&self, oid: ObjectId, data: UpdateSpaceshipDTO) -> Result<()>;
}

impl MongoSpaceshipRepository {
    pub fn new(collection: Collection<Spaceship>) -> Self {
        Self { collection }
    }
}

#[async_trait]
impl SpaceshipRepository for MongoSpaceshipRepository {
    async fn create(&self, data: CreateSpaceshipDTO) -> Result<ObjectId> {
        let id = self
            .collection
            .insert_one(Spaceship {
                _id: ObjectId::new(),
                name: data.name,
                user_id: data.user_id,
                location_id: data.location_id,
                flown_out_at: None,
                flying: false,
                system_id: data.system_id,
                planet_id: None,
            })
            .await
            .map_err(|_| CoreError::ServerError)?
            .inserted_id
            .as_object_id()
            .ok_or(CoreError::ServerError)?;

        Ok(id)
    }

    async fn find_one(&self, oid: ObjectId) -> Result<Option<Spaceship>> {
        let spaceship = self
            .collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|e| {
                tracing::error!("{:?}", e);
                CoreError::ServerError
            })?;

        Ok(spaceship)
    }

    async fn find_all(&self, filter: Document) -> Result<Vec<Spaceship>> {
        let cursor = self
            .collection
            .find(filter)
            .await
            .map_err(|_| CoreError::ServerError)?;

        cursor
            .try_collect()
            .await
            .map_err(|_| CoreError::ServerError)
    }

    async fn delete(&self, oid: ObjectId) -> Result<()> {
        self.collection
            .delete_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(())
    }

    async fn update(&self, oid: ObjectId, data: UpdateSpaceshipDTO) -> Result<()> {
        let mut update = Document::new();

        let serialized_data = bson::to_document(&data).map_err(|_| CoreError::ServerError)?;
        for (key, value) in serialized_data {
            let val = match key.as_str() {
                "flown_out_at" => match value {
                    Bson::DateTime(d) => Bson::Int64(d.to_chrono().timestamp()),
                    _ => value,
                },
                _ => value,
            };
            update.insert(key, val);
        }

        self.collection
            .find_one_and_update(doc! {"_id": oid}, doc! {"$set": update})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(())
    }
}
