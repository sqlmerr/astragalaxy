use mongodb::bson::{oid::ObjectId, Document};

use crate::{
    errors::CoreError,
    repositories::planet::{CreatePlanetDTO, PlanetRepository, UpdatePlanetDTO},
    schemas::planet::{CreatePlanetSchema, PlanetSchema, UpdatePlanetSchema},
    Result,
};

#[derive(Clone)]
pub struct PlanetService<R: PlanetRepository> {
    repository: R,
}

impl<R: PlanetRepository> PlanetService<R> {
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn create_planet(&self, data: CreatePlanetSchema) -> Result<PlanetSchema> {
        let dto = CreatePlanetDTO {
            system_id: data.system_id,
            threat: data.threat,
        };
        let planet_id = self.repository.create(dto).await?;
        self.find_one_planet(planet_id).await
    }

    pub async fn find_one_planet(&self, oid: ObjectId) -> Result<PlanetSchema> {
        let planet = self.repository.find_one(oid).await?;
        match planet {
            None => Err(CoreError::NotFound),
            Some(planet) => Ok(planet.into()),
        }
    }

    pub async fn find_all_planets(&self, filter: Document) -> Result<Vec<PlanetSchema>> {
        let planets: Vec<PlanetSchema> = self
            .repository
            .find_all(filter)
            .await?
            .iter()
            .map(|v| PlanetSchema::from(v.clone()))
            .collect();

        Ok(planets)
    }

    pub async fn delete_planet(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }

    pub async fn update_planet(&self, oid: ObjectId, data: UpdatePlanetSchema) -> Result<()> {
        let dto = UpdatePlanetDTO {
            system_id: data.system_id,
            threat: data.threat,
        };

        self.repository.update(oid, dto).await
    }
}
