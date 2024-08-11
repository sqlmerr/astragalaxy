use mongodb::{
    bson::{doc, oid::ObjectId, Bson, Document},
    Collection,
};

use crate::{errors::CoreError, models::User, Result};

#[derive(Clone)]
pub struct UserRepository {
    collection: Collection<User>,
}

pub struct CreateUserDTO {
    pub username: String,
    pub address: String,
}

pub struct UpdateUserDTO {
    pub username: Option<String>,
}

impl UserRepository {
    pub fn new(collection: Collection<User>) -> Self {
        Self { collection }
    }

    pub async fn create(&self, data: CreateUserDTO) -> Result<ObjectId> {
        let user_id = self
            .collection
            .insert_one(User {
                username: data.username,
                ton_address: data.address,
                ..Default::default()
            })
            .await
            .map_err(|_| CoreError::UsernameAlreadyOccupied)?
            .inserted_id;

        match user_id.as_object_id() {
            Some(oid) => Ok(oid),
            None => Err(CoreError::CantCreateUser),
        }
    }

    pub async fn find_one(&self, oid: ObjectId) -> Result<Option<User>> {
        self.collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)
    }

    pub async fn find_one_by_username(&self, username: String) -> Result<Option<User>> {
        self.collection
            .find_one(doc! {"username": username})
            .await
            .map_err(|_| CoreError::ServerError)
    }

    pub async fn find_one_filters(&self, filters: Document) -> Result<Option<User>> {
        self.collection
            .find_one(filters)
            .await
            .map_err(|_| CoreError::ServerError)
    }

    pub async fn find_all(&self) -> Vec<User> {
        todo!()
    }

    pub async fn delete(&self, oid: ObjectId) -> Result<()> {
        self.collection
            .find_one_and_delete(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(())
    }

    pub async fn update(&self, oid: ObjectId, data: UpdateUserDTO) -> Result<()> {
        let mut update = Document::new();

        if let Some(username) = data.username {
            update.insert("username", Bson::String(username));
        }

        self.collection
            .find_one_and_update(doc! {"_id": oid}, doc! {"$set": update})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(())
    }
}
