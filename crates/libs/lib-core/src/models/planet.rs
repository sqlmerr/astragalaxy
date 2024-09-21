use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Deserialize, Serialize)]
pub enum PlanetThreat {
    #[serde(rename = "radiation")]
    Radiation,
    #[serde(rename = "toxins")]
    Toxins,
    #[serde(rename = "freezing")]
    Freezing,
    #[serde(rename = "heat")]
    Heat,
}

impl ToString for PlanetThreat {
    fn to_string(&self) -> String {
        match self {
            PlanetThreat::Radiation => "radiation",
            PlanetThreat::Toxins => "toxins",
            PlanetThreat::Freezing => "freezing",
            PlanetThreat::Heat => "heat",
        }
        .to_string()
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
pub struct Planet {
    pub _id: ObjectId,
    pub system_id: ObjectId,
    pub threat: PlanetThreat,
}
