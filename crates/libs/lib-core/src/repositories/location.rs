use mongodb::{
    bson::{doc, oid::ObjectId},
    Collection,
};

use crate::{errors::CoreError, models::location::Location, Result};

pub struct LocationRepository {
    collection: Collection<Location>,
}

pub struct CreateLocationDTO {
    pub code: String,
    pub multiplayer: bool,
}

impl LocationRepository {
    pub fn new(collection: Collection<Location>) -> Self {
        Self { collection }
    }

    pub async fn create(&self, data: CreateLocationDTO) -> Result<ObjectId> {
        let id = self
            .collection
            .insert_one(Location {
                code: data.code,
                multiplayer: data.multiplayer,
                ..Default::default()
            })
            .await
            .map_err(|_| CoreError::ServerError)?
            .inserted_id
            .as_object_id()
            .ok_or(CoreError::ServerError)?;

        Ok(id)
    }

    pub async fn find_one(&self, oid: ObjectId) -> Result<Option<Location>> {
        let location = self
            .collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(location)
    }

    pub async fn find_one_by_code(&self, code: String) -> Result<Option<Location>> {
        let location = self
            .collection
            .find_one(doc! {"code": code})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(location)
    }

    pub async fn find_all(&self) -> Vec<Location> {
        todo!()
    }

    pub async fn delete(&self, oid: ObjectId) -> Result<()> {
        self.collection
            .delete_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;
        Ok(())
    }
}
