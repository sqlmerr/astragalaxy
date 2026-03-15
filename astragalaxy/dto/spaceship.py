from dataclasses import dataclass
from uuid import UUID

from astragalaxy.database.models import Spaceship


@dataclass(frozen=True)
class SpaceshipDTO:
    id: UUID
    name: str
    character_id: UUID
    active: bool
    point_id: str

    @classmethod
    def from_model(cls, model: Spaceship) -> "SpaceshipDTO":
        return cls(
            id=model.id,
            name=model.name,
            character_id=model.character_id,
            active=model.active,
            point_id=model.point_id,
        )


@dataclass(frozen=True)
class CreateSpaceshipDTO:
    name: str


@dataclass(frozen=True)
class RenameSpaceshipDTO:
    name: str
    spaceship_id: UUID
