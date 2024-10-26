use async_trait::async_trait;
use futures::stream::TryStreamExt;
use mongodb::{
    bson::{doc, oid::ObjectId, Bson, Document},
    Collection,
};

use crate::{errors::CoreError, models::System, Result};

#[derive(Clone)]
pub struct MongoSystemRepository {
    collection: Collection<System>,
}

pub struct CreateSystemDTO {
    pub name: String,
    pub neighbours: Vec<ObjectId>,
}

pub struct UpdateSystemDTO {
    pub name: Option<String>,
    pub neighbours: Option<Vec<ObjectId>>,
}

#[async_trait]
pub trait SystemRepository {
    async fn create(&self, data: CreateSystemDTO) -> Result<ObjectId>;
    async fn find_one(&self, oid: ObjectId) -> Result<Option<System>>;
    async fn find_one_by_name(&self, name: String) -> Result<Option<System>>;
    async fn find_all(&self, filter: Document) -> Result<Vec<System>>;
    async fn delete(&self, oid: ObjectId) -> Result<()>;
    async fn update(&self, oid: ObjectId, data: UpdateSystemDTO) -> Result<()>;
}

impl MongoSystemRepository {
    pub fn new(collection: Collection<System>) -> Self {
        Self { collection }
    }
}

#[async_trait]
impl SystemRepository for MongoSystemRepository {
    async fn create(&self, data: CreateSystemDTO) -> Result<ObjectId> {
        let id = self
            .collection
            .insert_one(System {
                name: data.name,
                neighbours: data.neighbours,
                ..Default::default()
            })
            .await
            .map_err(|_| CoreError::ServerError)?
            .inserted_id
            .as_object_id()
            .ok_or(CoreError::ServerError)?;

        Ok(id)
    }

    async fn find_one(&self, oid: ObjectId) -> Result<Option<System>> {
        let system = self
            .collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(system)
    }

    async fn find_one_by_name(&self, name: String) -> Result<Option<System>> {
        let system = self
            .collection
            .find_one(doc! {"name": name})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(system)
    }

    async fn find_all(&self, filter: Document) -> Result<Vec<System>> {
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
            .find_one_and_delete(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;
        Ok(())
    }

    async fn update(&self, oid: ObjectId, data: UpdateSystemDTO) -> Result<()> {
        let mut update = Document::new();

        if let Some(name) = data.name {
            update.insert("name", Bson::String(name));
        }

        if let Some(neighbours) = data.neighbours {
            update.insert(
                "neighbours",
                Bson::Array(neighbours.iter().map(|v| Bson::ObjectId(*v)).collect()),
            );
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
