pub mod errors;
pub mod jwt;
pub mod password;
pub mod schemas;

#[cfg(feature = "utoipa")]
use utoipa::openapi::{Object, ObjectBuilder};

#[cfg(feature = "utoipa")]
pub(crate) fn object_id_openapi() -> Object {
    ObjectBuilder::new()
        .schema_type(utoipa::openapi::SchemaType::String)
        .description(Some("ObjectId"))
        .build()
}
