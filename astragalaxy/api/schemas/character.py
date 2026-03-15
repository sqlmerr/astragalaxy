from uuid import UUID

from pydantic import BaseModel

from astragalaxy.dto.character import CharacterDTO


class CreateCharacterSchema(BaseModel):
    code: str


class CharacterSchema(BaseModel):
    id: UUID
    code: str
    user_id: UUID
    in_spaceship: bool
    point_id: str

    @classmethod
    def from_dto(cls, dto: CharacterDTO) -> "CharacterSchema":
        return cls.model_validate(dto, from_attributes=True)
