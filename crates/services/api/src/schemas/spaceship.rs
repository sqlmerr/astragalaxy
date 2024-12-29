use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
pub struct SpaceshipRenameSchema {
    pub name: String,
}
