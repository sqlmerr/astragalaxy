use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct OkResponse {
    ok: bool,
    custom_status_code: i32,
}

impl OkResponse {
    pub fn new(ok: bool) -> Self {
        Self {
            ok,
            custom_status_code: 1,
        }
    }

    pub fn status(self, custom_status_code: i32) -> Self {
        Self {
            ok: self.ok,
            custom_status_code: custom_status_code,
        }
    }
}
