use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::spaceship::{CreateSpaceshipDTO, SpaceshipRepository, UpdateSpaceshipDTO},
    schemas::spaceship::{CreateSpaceshipSchema, SpaceshipSchema, UpdateSpaceshipSchema},
    Result,
};

#[derive(Clone)]
pub struct SpaceshipService<R: SpaceshipRepository> {
    repository: R,
}

impl<R: SpaceshipRepository> SpaceshipService<R> {
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn create_spaceship(&self, data: CreateSpaceshipSchema) -> Result<SpaceshipSchema> {
        let dto = CreateSpaceshipDTO {
            name: data.name,
            user_id: data.user_id,
            location_id: data.location_id,
        };

        let spaceship_id = self.repository.create(dto).await?;
        match self.repository.find_one(spaceship_id).await? {
            Some(s) => Ok(s.into()),
            None => Err(CoreError::NotFound),
        }
    }

    pub async fn find_one_spaceship(&self, oid: ObjectId) -> Result<SpaceshipSchema> {
        match self.repository.find_one(oid).await? {
            None => Err(CoreError::NotFound),
            Some(spaceship) => Ok(spaceship.into()),
        }
    }

    pub async fn delete_spaceship(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }

    pub async fn update_spaceship(&self, oid: ObjectId, data: UpdateSpaceshipSchema) -> Result<()> {
        let dto = UpdateSpaceshipDTO {
            name: data.name,
            user_id: data.user_id,
            location_id: data.location_id,
        };

        self.repository.update(oid, dto).await
    }
}
