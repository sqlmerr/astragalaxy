from uuid import UUID

from pydantic import BaseModel

from astragalaxy.dto.spaceship import SpaceshipDTO


class SpaceshipSchema(BaseModel):
    id: UUID
    name: str
    character_id: UUID
    active: bool
    point_id: str

    @classmethod
    def from_dto(cls, dto: SpaceshipDTO) -> "SpaceshipSchema":
        return cls.model_validate(dto, from_attributes=True)


class CreateSpaceshipSchema(BaseModel):
    name: str


class RenameSpaceshipSchema(BaseModel):
    name: str
