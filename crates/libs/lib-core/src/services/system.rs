use mongodb::bson::oid::ObjectId;

use crate::{
    errors::CoreError,
    repositories::system::{CreateSystemDTO, SystemRepository, UpdateSystemDTO},
    schemas::system::{CreateSystemSchema, SystemSchema, UpdateSystemSchema},
    Result,
};

#[derive(Clone)]
pub struct SystemService<R: SystemRepository> {
    repository: R,
}

impl<R: SystemRepository> SystemService<R> {
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub async fn create_system(&self, data: CreateSystemSchema) -> Result<SystemSchema> {
        let dto = CreateSystemDTO {
            name: data.name,
            neighbours: data.neighbours.clone(),
        };
        let system_id = self.repository.create(dto).await?;
        for neighbour in data.neighbours {
            let other_system = self.find_one_system(neighbour).await?;
            let mut n = other_system.neighbours.clone();
            n.push(system_id);
            self.repository
                .update(
                    neighbour,
                    UpdateSystemDTO {
                        neighbours: Some(n),
                        name: None,
                    },
                )
                .await?;
        }

        let system = self.find_one_system(system_id).await?;

        Ok(system)
    }

    pub async fn find_one_system(&self, oid: ObjectId) -> Result<SystemSchema> {
        let system = self.repository.find_one(oid).await?;
        match system {
            Some(s) => Ok(s.into()),
            None => Err(CoreError::NotFound),
        }
    }

    pub async fn delete_system(&self, oid: ObjectId) -> Result<()> {
        self.repository.delete(oid).await
    }

    pub async fn update_system(&self, oid: ObjectId, data: UpdateSystemSchema) -> Result<()> {
        if let Some(neighbours) = data.neighbours.clone() {
            for neighbour in neighbours {
                self.find_one_system(neighbour).await?;
            }
        }

        let dto = UpdateSystemDTO {
            name: data.name,
            neighbours: data.neighbours,
        };

        self.repository.update(oid, dto).await
    }

    pub async fn find_one_system_by_name(&self, name: String) -> Result<SystemSchema> {
        match self.repository.find_one_by_name(name).await? {
            Some(s) => Ok(s.into()),
            None => Err(CoreError::NotFound),
        }
    }
}
