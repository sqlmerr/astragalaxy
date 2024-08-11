use mongodb::bson::doc;
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

    pub async fn register(&self, data: CreateUserSchema, address: String) -> Result<UserSchema> {
        if let Some(_) = self
            .repository
            .find_one_by_username(data.username.clone())
            .await?
        {
            return Err(CoreError::UsernameAlreadyOccupied);
        }

        if let Some(_) = self
            .repository
            .find_one_filters(doc! {"ton_address": address.clone()})
            .await?
        {
            return Err(CoreError::AddressAlreadyOccupied);
        }

        let dto = CreateUserDTO {
            username: data.username,
            address,
        };
        let user_id = self.repository.create(dto).await?;
        let user = self.repository.find_one(user_id).await?;

        match user {
            Some(user) => Ok(UserSchema::from(user)),
            None => Err(CoreError::CantCreateUser),
        }
    }

    pub async fn login(&self, _username: String, _password: String) -> Result<String> {
        todo!()
    }

    pub async fn find_one_user(&self, oid: ObjectId) -> Result<UserSchema> {
        let user = self.repository.find_one(oid).await?;

        match user {
            Some(u) => Ok(UserSchema::from(u)),
            None => Err(CoreError::NotFound),
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

    pub async fn find_one_user_by_username(&self, username: String) -> Result<UserSchema> {
        let user = self.repository.find_one_by_username(username).await?;
        if let Some(u) = user {
            return Ok(u.into());
        }
        Err(CoreError::NotFound)
    }

    pub async fn find_one_user_by_address(&self, address: String) -> Result<UserSchema> {
        let user = self
            .repository
            .find_one_filters(doc! {"ton_address": address})
            .await?;
        if let Some(u) = user {
            return Ok(u.into());
        }
        Err(CoreError::NotFound)
    }

    pub async fn delete_user(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }
}
