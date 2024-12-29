use async_trait::async_trait;
use chrono::{NaiveDateTime, TimeZone, Utc};
use mongodb::{
    bson::{doc, oid::ObjectId, Bson, DateTime, Document},
    Collection,
};

use crate::{errors::CoreError, models::Spaceship, Result};

#[derive(Clone)]
pub struct MongoSpaceshipRepository {
    collection: Collection<Spaceship>,
}

pub struct CreateSpaceshipDTO {
    pub name: String,
    pub user_id: ObjectId,
    pub location_id: ObjectId,
    pub system_id: ObjectId,
}

#[derive(Default)]
pub struct UpdateSpaceshipDTO {
    pub name: Option<String>,
    pub user_id: Option<ObjectId>,
    pub location_id: Option<ObjectId>,
    pub flown_out_at: Option<Option<NaiveDateTime>>,
    pub flying: Option<bool>,
    pub system_id: Option<ObjectId>,
    pub planet_id: Option<Option<ObjectId>>,
}

#[async_trait]
pub trait SpaceshipRepository {
    async fn create(&self, data: CreateSpaceshipDTO) -> Result<ObjectId>;
    async fn find_one(&self, oid: ObjectId) -> Result<Option<Spaceship>>;
    async fn delete(&self, oid: ObjectId) -> Result<()>;
    async fn update(&self, oid: ObjectId, data: UpdateSpaceshipDTO) -> Result<()>;
}

impl MongoSpaceshipRepository {
    pub fn new(collection: Collection<Spaceship>) -> Self {
        Self { collection }
    }
}

#[async_trait]
impl SpaceshipRepository for MongoSpaceshipRepository {
    async fn create(&self, data: CreateSpaceshipDTO) -> Result<ObjectId> {
        let id = self
            .collection
            .insert_one(Spaceship {
                _id: ObjectId::new(),
                name: data.name,
                user_id: data.user_id,
                location_id: data.location_id,
                flown_out_at: None,
                flying: false,
                system_id: data.system_id,
                planet_id: None,
            })
            .await
            .map_err(|_| CoreError::ServerError)?
            .inserted_id
            .as_object_id()
            .ok_or(CoreError::ServerError)?;

        Ok(id)
    }

    async fn find_one(&self, oid: ObjectId) -> Result<Option<Spaceship>> {
        let spaceship = self
            .collection
            .find_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(spaceship)
    }

    async fn delete(&self, oid: ObjectId) -> Result<()> {
        self.collection
            .delete_one(doc! {"_id": oid})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(())
    }

    async fn update(&self, oid: ObjectId, data: UpdateSpaceshipDTO) -> Result<()> {
        let mut update = Document::new();

        if let Some(name) = data.name {
            update.insert("name", Bson::String(name));
        }

        if let Some(user_id) = data.user_id {
            update.insert("user_id", Bson::ObjectId(user_id));
        }

        if let Some(location_id) = data.location_id {
            update.insert("location_id", Bson::ObjectId(location_id));
        }

        if let Some(flown_out_at) = data.flown_out_at {
            match flown_out_at {
                None => update.insert("flown_out_at", Bson::Null),
                Some(f) => {
                    let utc_datetime = Utc.from_utc_datetime(&f);
                    update.insert(
                        "flown_out_at",
                        Bson::DateTime(DateTime::from_chrono(utc_datetime)),
                    )
                }
            };
        }

        if let Some(flying) = data.flying {
            update.insert("flying", Bson::Boolean(flying));
        }

        if let Some(system_id) = data.system_id {
            update.insert("system_id", Bson::ObjectId(system_id));
        }

        if let Some(planet_id) = data.planet_id {
            let val = match planet_id {
                None => Bson::Null,
                Some(id) => Bson::ObjectId(id),
            };
            update.insert("planet_id", val);
        }

        self.collection
            .find_one_and_update(doc! {"_id": oid}, doc! {"$set": update})
            .await
            .map_err(|_| CoreError::ServerError)?;

        Ok(())
    }
}
