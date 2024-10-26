use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct OkResponse {
    ok: bool,
}

impl OkResponse {
    pub fn new(ok: bool) -> Self {
        Self { ok }
    }
}
