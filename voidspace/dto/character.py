from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models import Character


@dataclass(frozen=True)
class CharacterDTO:
    id: UUID
    code: str
    location: str
    user_id: UUID

    @classmethod
    def from_model(cls, model: Character) -> "CharacterDTO":
        return cls(
            id=model.id, code=model.code, location=model.location, user_id=model.user_id
        )


@dataclass(frozen=True)
class CreateCharacterDTO:
    user_id: UUID
    code: str
