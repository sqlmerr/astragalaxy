from uuid import UUID

from pydantic import BaseModel

from astragalaxy.dto.spaceship import SpaceshipDTO


class SpaceshipSchema(BaseModel):
    id: UUID
    name: str
    location: str
    character_id: UUID
    active: bool
    system_id: str
    planet_id: str | None

    @classmethod
    def from_dto(cls, dto: SpaceshipDTO) -> "SpaceshipSchema":
        return cls(
            id=dto.id,
            name=dto.name,
            location=dto.location,
            character_id=dto.character_id,
            active=dto.active,
            system_id=dto.system_id,
            planet_id=dto.planet_id,
        )


class CreateSpaceshipSchema(BaseModel):
    name: str


class RenameSpaceshipSchema(BaseModel):
    name: str
