use async_trait::async_trait;
use futures::stream::TryStreamExt;
use mongodb::{
    bson::{doc, oid::ObjectId, Document},
    Collection,
};
use serde::Serialize;

use crate::{
    errors::CoreError,
    models::{planet::PlanetThreat, Planet},
    Result,
};

#[derive(Clone)]
pub struct MongoPlanetRepository {
    collection: Collection<Planet>,
}

pub struct CreatePlanetDTO {
    pub system_id: ObjectId,
    pub threat: PlanetThreat,
}

#[derive(Serialize)]
pub struct UpdatePlanetDTO {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub system_id: Option<ObjectId>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub threat: Option<PlanetThreat>,
}

#[async_trait]
pub trait PlanetRepository {
    async fn create(&self, data: CreatePlanetDTO) -> Result<ObjectId>;
    async fn find_one(&self, oid: ObjectId) -> Result<Option<Planet>>;
    async fn find_all(&self, filter: Document) -> Result<Vec<Planet>>;
    async fn delete(&self, oid: ObjectId) -> Result<()>;
    async fn update(&self, oid: ObjectId, data: UpdatePlanetDTO) -> Result<()>;
}

impl MongoPlanetRepository {
    pub fn new(collection: Collection<Planet>) -> Self {
        Self { collection }
    }
}

#[async_trait]
impl PlanetRepository for MongoPlanetRepository {
    async fn create(&self, data: CreatePlanetDTO) -> Result<ObjectId> {
        let id = self
            .collection
            .insert_one(Planet {
                _id: ObjectId::new(),
                system_id: data.system_id,
                threat: data.threat,
            })
            .await
            .map_err(|_| CoreError::ServerError)?
            .inserted_id
            .as_object_id()
            .ok_or(CoreError::ServerError)?;

        Ok(id)
    }

    async fn find_one(&self, oid: ObjectId) -> Result<Option<Planet>> {
        let planet = self
            .collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(planet)
    }

    async fn find_all(&self, filter: Document) -> Result<Vec<Planet>> {
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

    async fn update(&self, oid: ObjectId, data: UpdatePlanetDTO) -> Result<()> {
        let mut update = Document::new();

        let serialized_data = bson::to_document(&data).map_err(|_| CoreError::ServerError)?;
        for (key, value) in serialized_data {
            let val = match key.as_str() {
                _ => value,
            };
            update.insert(key, val);
        }

        self.collection
            .find_one_and_update(doc! {"_id": oid}, doc! {"$set": update})
            .await
            .map_err(|e| {
                eprintln!("{}", e);
                CoreError::ServerError
            })?;

        Ok(())
    }
}
