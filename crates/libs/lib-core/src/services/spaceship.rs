use bson::Document;
use chrono::{TimeDelta, Utc};
use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::{
        planet::PlanetRepository,
        spaceship::{CreateSpaceshipDTO, SpaceshipRepository, UpdateSpaceshipDTO},
        system::SystemRepository,
    },
    schemas::spaceship::{CreateSpaceshipSchema, SpaceshipSchema, UpdateSpaceshipSchema},
    Result,
};

use super::{planet::PlanetService, system::SystemService};

#[derive(Clone)]
pub struct SpaceshipService<R: SpaceshipRepository, P: PlanetRepository, S: SystemRepository> {
    repository: R,
    planet_service: PlanetService<P>,
    system_service: SystemService<S>,
}

impl<R: SpaceshipRepository, P: PlanetRepository, S: SystemRepository> SpaceshipService<R, P, S> {
    pub fn new(
        repository: R,
        planet_service: PlanetService<P>,
        system_service: SystemService<S>,
    ) -> Self {
        Self {
            repository,
            planet_service,
            system_service,
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

    pub async fn find_all_spaceships(&self, filter: Document) -> Result<Vec<SpaceshipSchema>> {
        let spaceships = self.repository.find_all(filter).await?;
        Ok(spaceships
            .iter()
            .map(|v| SpaceshipSchema::from(v.clone()))
            .collect())
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

        if spaceship.flown_out_at.is_some() {
            let now = Utc::now().timestamp();
            let difference = now - spaceship.flown_out_at.unwrap();
            if difference >= TimeDelta::minutes(10).num_seconds() {
                let dto = UpdateSpaceshipDTO {
                    flying: Some(false),
                    flown_out_at: Some(None),
                    ..Default::default()
                };
                self.repository.update(oid, dto).await?;
            } else {
                return Err(CoreError::SpaceshipAlreadyFlying);
            }
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
            flown_out_at: Some(Some(Utc::now())),
            planet_id: Some(Some(planet._id)),
            ..Default::default()
        };
        self.repository.update(oid, dto).await?;

        Ok(())
    }

    pub async fn hyperjump(&self, oid: ObjectId, system_id: ObjectId) -> Result<()> {
        let spaceship = match self.repository.find_one(oid).await? {
            None => return Err(CoreError::PlayerHasNoSpaceship),
            Some(s) => s,
        };

        if spaceship.flying && spaceship.flown_out_at.is_some() {
            return Err(CoreError::SpaceshipAlreadyFlying);
        } else if spaceship.flying && spaceship.flown_out_at.is_none() {
            return Err(CoreError::ServerError);
        }

        let system = self.system_service.find_one_system(system_id).await?;
        if system._id == spaceship.system_id {
            return Err(CoreError::SpaceshipIsAlreadyInThisSystem);
        }

        let dto = UpdateSpaceshipDTO {
            flown_out_at: Some(Some(Utc::now())),
            system_id: Some(system._id),
            ..Default::default()
        };
        self.repository.update(oid, dto).await?;

        Ok(())
    }
}
