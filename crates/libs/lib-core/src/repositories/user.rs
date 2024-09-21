use async_trait::async_trait;
use mongodb::{
    bson::{doc, oid::ObjectId, Bson, Document},
    Collection,
};

use crate::{errors::CoreError, models::User, Result};

#[derive(Clone)]
pub struct MongoUserRepository {
    collection: Collection<User>,
}

#[derive(Default)]
pub struct CreateUserDTO {
    pub username: String,
    pub discord_id: Option<i64>,
    pub password: Option<String>,
    pub location_id: ObjectId,
}

#[derive(Default)]
pub struct UpdateUserDTO {
    pub username: Option<String>,
    pub password: Option<String>,
    pub x: Option<i64>,
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
                password: data.password,
                discord_id: data.discord_id,
                location_id: data.location_id,
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

        if let Some(username) = data.username {
            update.insert("username", Bson::String(username));
        }

        if let Some(password) = data.password {
            update.insert("password", Bson::String(password));
        }

        if let Some(x) = data.x {
            update.insert("x", Bson::Int64(x));
        }

        if let Some(y) = data.y {
            update.insert("y", Bson::Int64(y));
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
