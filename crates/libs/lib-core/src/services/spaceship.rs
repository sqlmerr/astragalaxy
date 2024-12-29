use chrono::Utc;
use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::{
        planet::PlanetRepository,
        spaceship::{CreateSpaceshipDTO, SpaceshipRepository, UpdateSpaceshipDTO},
    },
    schemas::spaceship::{CreateSpaceshipSchema, SpaceshipSchema, UpdateSpaceshipSchema},
    Result,
};

use super::planet::PlanetService;

#[derive(Clone)]
pub struct SpaceshipService<R: SpaceshipRepository, P: PlanetRepository> {
    repository: R,
    planet_service: PlanetService<P>,
}

impl<R: SpaceshipRepository, P: PlanetRepository> SpaceshipService<R, P> {
    pub fn new(repository: R, planet_service: PlanetService<P>) -> Self {
        Self {
            repository,
            planet_service,
        }
    }

    pub async fn create_spaceship(&self, data: CreateSpaceshipSchema) -> Result<SpaceshipSchema> {
        let dto = CreateSpaceshipDTO {
            name: data.name,
            user_id: data.user_id,
            location_id: data.location_id,
            system_id: data.system_id,
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
            flown_out_at: data.flown_out_at,
            flying: data.flying,
            system_id: data.system_id,
            planet_id: data.planet_id,
        };

        self.repository.update(oid, dto).await
    }

    pub async fn fly(&self, oid: ObjectId, planet_id: ObjectId) -> Result<()> {
        let spaceship = match self.repository.find_one(oid).await? {
            None => return Err(CoreError::PlayerHasNoSpaceship),
            Some(s) => s,
        };
        if spaceship.flying && spaceship.flown_out_at.is_some() {
            return Err(CoreError::SpaceshipAlreadyFlying);
        } else if spaceship.flying && spaceship.flown_out_at.is_none() {
            return Err(CoreError::ServerError);
        }

        let planet = self.planet_service.find_one_planet(planet_id).await?;
        if planet.system_id != spaceship.system_id {
            return Err(CoreError::PlanetIsInAnotherSystem);
        }

        if spaceship.planet_id.is_some() && planet._id == spaceship.planet_id.unwrap() {
            return Err(CoreError::SpaceshipIsAlreadyInThisPlanet);
        }

        let dto = UpdateSpaceshipDTO {
            flying: Some(true),
            flown_out_at: Some(Some(Utc::now().naive_utc())),
            planet_id: Some(Some(planet._id)),
            ..Default::default()
        };
        self.repository.update(oid, dto).await?;

        Ok(())
    }
}
