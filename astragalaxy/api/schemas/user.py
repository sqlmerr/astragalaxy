from uuid import UUID

from pydantic import BaseModel

from astragalaxy.dto.user import UserDTO


class UserSchema(BaseModel):
    id: UUID
    username: str

    @classmethod
    def from_dto(cls, dto: UserDTO) -> "UserSchema":
        return cls(id=dto.id, username=dto.username)
