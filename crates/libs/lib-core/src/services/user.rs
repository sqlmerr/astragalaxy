use lib_auth::{
    errors::AuthError,
    jwt::create_token,
    password::{hash_password, verify_password},
    schemas::Claims,
};
use mongodb::bson::{doc, oid::ObjectId};

use crate::{
    errors::CoreError,
    repositories::user::{CreateUserDTO, UpdateUserDTO, UserRepository},
    schemas::user::{CreateUserSchema, UpdateUserSchema, UserSchema},
    Result,
};

#[derive(Clone)]
pub struct UserService<R: UserRepository> {
    repository: R,
}

impl<R: UserRepository> UserService<R> {
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn register(
        &self,
        data: CreateUserSchema,
        location_id: ObjectId,
    ) -> Result<UserSchema> {
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
            password: Some(hashed_password),
            location_id,
            ..Default::default()
        };
        let user_id = self.repository.create(dto).await?;
        let user = self.repository.find_one(user_id).await?;

        match user {
            Some(user) => Ok(UserSchema::from(user)),
            None => Err(CoreError::CantCreateUser),
        }
    }

    pub async fn register_from_discord(
        &self,
        discord_id: i64,
        username: String,
        location_id: ObjectId,
    ) -> Result<UserSchema> {
        if let Some(_) = self
            .repository
            .find_one_by_username(username.clone())
            .await?
        {
            return Err(CoreError::UsernameAlreadyOccupied);
        }

        let dto = CreateUserDTO {
            username,
            password: None,
            discord_id: Some(discord_id),
            location_id,
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
                if u.password.is_none() {
                    return Err(AuthError::WrongCredentials.into());
                }

                if !verify_password(password, u.password.unwrap()) {
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
            None => Err(CoreError::NotFound),
        }
    }

    pub async fn find_one_user_by_discord_id(&self, discord_id: i64) -> Result<UserSchema> {
        let user = self
            .repository
            .find_one_filters(doc! { "discord_id": discord_id })
            .await?;
        match user {
            Some(u) => Ok(UserSchema::from(u)),
            None => Err(CoreError::NotFound),
        }
    }

    // pub async fn find_all_users(&self, filters: Doc) -> Vec<UserSchema> {
    //     self.repository
    //         .find_all()
    //         .await
    //         .iter()
    //         .map(|v| UserSchema::from(v.clone()))
    //         .collect()
    // }

    pub async fn update_user(&self, oid: ObjectId, data: UpdateUserSchema) -> Result<()> {
        let hashed_password;
        if let Some(password) = data.password {
            hashed_password = Some(hash_password(password))
        } else {
            hashed_password = None
        }

        let dto = UpdateUserDTO {
            username: data.username,
            password: hashed_password,
            spaceship_id: data.spaceship_id,
            ..Default::default()
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

    pub async fn delete_user(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }

    pub async fn move_user(&self, oid: ObjectId, x: i64, y: i64) -> Result<()> {
        let dto = UpdateUserDTO {
            x: Some(x),
            y: Some(y),
            ..Default::default()
        };
        self.repository.update(oid, dto).await
    }

    pub async fn get_users_count_by_location(&self, location_id: ObjectId) -> u64 {
        match self
            .repository
            .get_count_filters(doc! { "location_id": location_id })
            .await
        {
            Ok(count) => count,
            Err(_) => 0,
        }
    }
}
