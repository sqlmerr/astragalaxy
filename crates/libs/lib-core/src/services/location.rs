use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::location::{CreateLocationDTO, LocationRepository},
    schemas::location::{CreateLocationSchema, LocationSchema},
    Result,
};

#[derive(Clone)]
pub struct LocationService<R: LocationRepository> {
    repository: R,
}

impl<R> LocationService<R>
where
    R: LocationRepository,
{
    pub fn new(repository: R) -> Self {
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

    pub async fn find_one_location_by_code(&self, code: String) -> Result<LocationSchema> {
        let location = self.repository.find_one_by_code(code).await?;
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

#[cfg(test)]
mod tests {
    use async_trait::async_trait;
    use mongodb::bson::oid::ObjectId;

    use crate::{
        models::location::Location,
        repositories::location::{CreateLocationDTO, LocationRepository},
        schemas::location::CreateLocationSchema,
        services::location::LocationService,
        Result,
    };

    #[derive(Clone, Default)]
    pub struct MockLocationRepository;

    #[async_trait]
    impl LocationRepository for MockLocationRepository {
        async fn create(&self, data: CreateLocationDTO) -> Result<ObjectId> {
            let location = Location {
                _id: ObjectId::new(),
                code: data.code,
                multiplayer: data.multiplayer,
            };

            Ok(location._id)
        }
        async fn find_one(&self, oid: ObjectId) -> Result<Option<Location>> {
            Ok(Some(Location {
                _id: oid,
                code: "test".to_string(),
                multiplayer: false,
            }))
        }
        async fn find_one_by_code(&self, code: String) -> Result<Option<Location>> {
            Ok(Some(Location {
                _id: ObjectId::new(),
                code,
                multiplayer: false,
            }))
        }
        async fn find_all(&self) -> Vec<Location> {
            vec![Location {
                _id: ObjectId::new(),
                code: "test".to_string(),
                multiplayer: false,
            }]
        }
        async fn delete(&self, _oid: ObjectId) -> Result<()> {
            Ok(())
        }
    }

    #[tokio::test]
    async fn test_create_location() {
        let repository = MockLocationRepository::default();
        let service = LocationService::new(repository);
        let data = CreateLocationSchema {
            code: "test".to_string(),
            multiplayer: false,
        };
        let loc = service.create_location(data).await.unwrap();
        assert_eq!(loc.code, "test");
        assert_eq!(loc.multiplayer, false);
    }

    #[tokio::test]
    async fn test_find_one_location() {
        let repository = MockLocationRepository::default();
        let service = LocationService::new(repository);
        let id = ObjectId::new();
        let loc = service.find_one_location(id).await.unwrap();

        assert_eq!(loc._id, id);
        assert_eq!(loc.code, "test");
        assert_eq!(loc.multiplayer, false);
    }

    #[tokio::test]
    async fn test_find_one_location_by_code() {
        let repository = MockLocationRepository::default();
        let service = LocationService::new(repository);
        let code = "test".to_string();
        let loc = service
            .find_one_location_by_code(code.clone())
            .await
            .unwrap();
        assert_eq!(loc.code, code);
        assert_eq!(loc.multiplayer, false);
    }

    #[tokio::test]
    async fn test_delete_location() {
        let repository = MockLocationRepository::default();
        let service = LocationService::new(repository);
        let id = ObjectId::new();
        let response = service.delete_location(id).await.unwrap();

        assert_eq!(response, ())
    }
}
