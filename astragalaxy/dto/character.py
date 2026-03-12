from dataclasses import dataclass
from uuid import UUID

from astragalaxy.database.models import Character


@dataclass(frozen=True)
class CharacterDTO:
    id: UUID
    code: str
    location: str
    in_spaceship: bool
    user_id: UUID
    system_id: str

    @classmethod
    def from_model(cls, model: Character) -> "CharacterDTO":
        return cls(
            id=model.id,
            code=model.code,
            location=model.location,
            in_spaceship=model.in_spaceship,
            user_id=model.user_id,
            system_id=model.system_id,
        )


@dataclass(frozen=True)
class CreateCharacterDTO:
    code: str
