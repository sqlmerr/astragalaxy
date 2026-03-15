from dataclasses import dataclass
from uuid import UUID

from astragalaxy.database.models import Character


@dataclass(frozen=True)
class CharacterDTO:
    id: UUID
    code: str
    in_spaceship: bool
    user_id: UUID
    point_id: str

    @classmethod
    def from_model(cls, model: Character) -> "CharacterDTO":
        return cls(
            id=model.id,
            code=model.code,
            in_spaceship=model.in_spaceship,
            user_id=model.user_id,
            point_id=model.point_id
        )


@dataclass(frozen=True)
class CreateCharacterDTO:
    code: str
