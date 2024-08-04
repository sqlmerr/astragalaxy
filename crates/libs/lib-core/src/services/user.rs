use lib_auth::{
    errors::AuthError,
    jwt::create_token,
    password::{hash_password, verify_password},
    schemas::Claims,
};
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

    pub async fn register(&self, data: CreateUserSchema) -> Result<UserSchema> {
        if let Some(_) = self
            .repository
            .find_one_by_username(data.username.clone())
            .await?
        {
            return Err(CoreError::UsernameAlreadyOccupied);
        }

        let hashed_password = hash_password(data.password);

        let dto = CreateUserDTO {
            username: data.username,
            hashed_password,
        };
        let user_id = self.repository.create(dto).await?;
        let user = self.repository.find_one(user_id).await?;

        match user {
            Some(user) => Ok(UserSchema::from(user)),
            None => Err(CoreError::CantCreateUser),
        }
    }

    pub async fn login(&self, username: String, password: String) -> Result<String> {
        let user = self
            .repository
            .find_one_by_username(username.clone())
            .await?;
        match user {
            None => Err(AuthError::WrongCredentials.into()),
            Some(u) => {
                if !verify_password(password, u.password) {
                    return Err(AuthError::WrongCredentials.into());
                }

                let claims = Claims::new(username);
                let token = create_token(&claims).map_err(|_| AuthError::TokenCreation)?;

                Ok(token)
            }
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
        let password;
        if let Some(p) = data.password {
            password = Some(hash_password(p));
        } else {
            password = None;
        }

        let dto = UpdateUserDTO {
            username: data.username,
            hashed_password: password,
        };
        self.repository.update(oid, dto).await
    }

    pub async fn find_one_user_by_username(&self, username: String) -> Result<UserSchema> {
        let user = self.repository.find_one_by_username(username).await?;
        if let Some(u) = user {
            return Ok(u.into());
        }
        Err(CoreError::UserNotFound)
    }

    pub async fn delete_user(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }
}
