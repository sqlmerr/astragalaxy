use async_trait::async_trait;
use lib_utils::generate_token;
use mongodb::{
    bson::{doc, oid::ObjectId, Document},
    Collection,
};
use serde::Serialize;

use crate::{errors::CoreError, models::User, Result};

#[derive(Clone)]
pub struct MongoUserRepository {
    collection: Collection<User>,
}

#[derive(Default)]
pub struct CreateUserDTO {
    pub username: String,
    pub discord_id: Option<i64>,
    pub telegram_id: i64,
    pub location_id: ObjectId,
    pub system_id: ObjectId,
}

#[derive(Default, Serialize)]
pub struct UpdateUserDTO {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub username: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub spaceship_id: Option<ObjectId>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub in_spaceship: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub system_id: Option<ObjectId>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub x: Option<i64>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub y: Option<i64>,
}

#[async_trait]
pub trait UserRepository {
    async fn create(&self, data: CreateUserDTO) -> Result<ObjectId>;
    async fn find_one(&self, oid: ObjectId) -> Result<Option<User>>;
    async fn find_one_by_username(&self, username: String) -> Result<Option<User>>;
    async fn find_one_filters(&self, filters: Document) -> Result<Option<User>>;
    async fn find_all(&self, filters: Document) -> Vec<User>;
    async fn get_count_filters(&self, filters: Document) -> Result<u64>;
    async fn delete(&self, oid: ObjectId) -> Result<()>;
    async fn update(&self, oid: ObjectId, data: UpdateUserDTO) -> Result<()>;
}

impl MongoUserRepository {
    pub fn new(collection: Collection<User>) -> Self {
        Self { collection }
    }
}

#[async_trait]
impl UserRepository for MongoUserRepository {
    async fn create(&self, data: CreateUserDTO) -> Result<ObjectId> {
        let user_id = self
            .collection
            .insert_one(User {
                username: data.username,
                discord_id: data.discord_id,
                telegram_id: data.telegram_id,
                location_id: data.location_id,
                system_id: data.system_id,
                in_spaceship: false,
                token: generate_token(32),
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

    async fn find_one(&self, oid: ObjectId) -> Result<Option<User>> {
        self.collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|e| {
                eprintln!("{}", e);
                CoreError::ServerError
            })
    }

    async fn find_one_by_username(&self, username: String) -> Result<Option<User>> {
        self.collection
            .find_one(doc! {"username": username})
            .await
            .map_err(|e| {
                eprintln!("{}", e);
                CoreError::ServerError
            })
    }

    async fn find_one_filters(&self, filters: Document) -> Result<Option<User>> {
        self.collection.find_one(filters).await.map_err(|e| {
            eprintln!("{}", e);
            CoreError::ServerError
        })
    }

    async fn find_all(&self, _filters: Document) -> Vec<User> {
        todo!()
    }

    async fn get_count_filters(&self, filters: Document) -> Result<u64> {
        self.collection.count_documents(filters).await.map_err(|e| {
            eprintln!("{}", e);
            CoreError::ServerError
        })
    }

    async fn delete(&self, oid: ObjectId) -> Result<()> {
        self.collection
            .find_one_and_delete(doc! {"_id": oid})
            .await
            .map_err(|e| {
                eprintln!("{}", e);
                CoreError::ServerError
            })?;

        Ok(())
    }

    async fn update(&self, oid: ObjectId, data: UpdateUserDTO) -> Result<()> {
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
