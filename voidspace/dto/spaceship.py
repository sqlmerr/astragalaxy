from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models import Spaceship


@dataclass(frozen=True)
class SpaceshipDTO:
    id: UUID
    name: str
    location: str
    character_id: UUID
    active: bool
    system_id: str
    planet_id: str | None

    @classmethod
    def from_model(cls, model: Spaceship) -> "SpaceshipDTO":
        return cls(
            id=model.id,
            name=model.name,
            location=model.location,
            character_id=model.character_id,
            active=model.active,
            system_id=model.system_id,
            planet_id=model.planet_id,
        )


@dataclass(frozen=True)
class CreateSpaceshipDTO:
    name: str
    system_id: str


@dataclass(frozen=True)
class RenameSpaceshipDTO:
    name: str
    spaceship_id: UUID
