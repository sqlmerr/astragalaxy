from uuid import UUID

from pydantic import BaseModel

from voidspace.dto.character import CharacterDTO


class CreateCharacterSchema(BaseModel):
    code: str


class CharacterSchema(BaseModel):
    id: UUID
    code: str
    location: str
    user_id: UUID
    system_id: str
    in_spaceship: bool

    @classmethod
    def from_dto(cls, dto: CharacterDTO) -> "CharacterSchema":
        return cls(
            id=dto.id,
            code=dto.code,
            location=dto.location,
            user_id=dto.user_id,
            system_id=dto.system_id,
            in_spaceship=dto.in_spaceship,
        )
