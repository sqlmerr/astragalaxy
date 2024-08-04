use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::user::{CreateUserDTO, UpdateUserDTO, UserRepository},
    schemas::user::{CreateUserSchema, UpdateUserSchema, UserSchema},
    Result,
};

#[derive(Clone)]
pub struct UserService {
    repository: UserRepository,
}

impl UserService {
    pub fn new(repository: UserRepository) -> Self {
        Self { repository }
    }

    pub async fn create_user(&self, data: CreateUserSchema) -> Result<UserSchema> {
        if let Some(u) = self
            .repository
            .fjnd_one_by_username(data.username.clone())
            .await?
        {
            return Err(CoreError::UsernameAlreadyOccupied);
        }

        let dto = CreateUserDTO {
            username: data.username,
        };
        let user_id = self.repository.create(dto).await?;
        let user = self.repository.find_one(user_id).await?;

        match user {
            Some(user) => Ok(UserSchema::from(user)),
            None => Err(CoreError::CantCreateUser),
        }
    }

    pub async fn find_one_user(&self, oid: ObjectId) -> Result<UserSchema> {
        let user = self.repository.find_one(oid).await?;

        match user {
            Some(u) => Ok(UserSchema::from(u)),
            None => Err(CoreError::UserNotFound),
        }
    }

    pub async fn find_all_users(&self) -> Vec<UserSchema> {
        self.repository
            .find_all()
            .await
            .iter()
            .map(|v| UserSchema::from(v.clone()))
            .collect()
    }

    pub async fn update_user(&self, oid: ObjectId, data: UpdateUserSchema) -> Result<()> {
        let dto = UpdateUserDTO {
            username: data.username,
        };
        self.repository.update(oid, dto).await
    }

    pub async fn delete_user(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }
}
