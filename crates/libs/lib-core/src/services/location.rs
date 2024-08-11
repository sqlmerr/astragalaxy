use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::location::{CreateLocationDTO, LocationRepository},
    schemas::location::{CreateLocationSchema, LocationSchema},
    Result,
};

pub struct LocationService {
    repository: LocationRepository,
}

impl LocationService {
    pub fn new(repository: LocationRepository) -> Self {
        Self { repository }
    }

    pub async fn create_location(&self, data: CreateLocationSchema) -> Result<LocationSchema> {
        let dto = CreateLocationDTO {
            code: data.code,
            multiplayer: data.multiplayer,
        };
        let loc_id = self.repository.create(dto).await?;
        let location = self.find_one_location(loc_id).await?;

        Ok(location)
    }

    pub async fn find_one_location(&self, oid: ObjectId) -> Result<LocationSchema> {
        let location = self.repository.find_one(oid).await?;
        match location {
            Some(loc) => Ok(loc.into()),
            None => Err(CoreError::NotFound),
        }
    }

    pub async fn find_all_locations(&self) -> Vec<LocationSchema> {
        todo!()
    }

    pub async fn delete_location(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }
}
