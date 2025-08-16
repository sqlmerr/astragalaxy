from dataclasses import dataclass
from uuid import UUID

from voidspace.database.models import User


@dataclass(frozen=True)
class UserDTO:
    id: UUID
    username: str
    password: str
    token: str

    @classmethod
    def from_model(cls, model: User) -> "UserDTO":
        return cls(
            id=model.id,
            username=model.username,
            password=model.password,
            token=model.token,
        )


@dataclass(frozen=True)
class CreateUserDTO:
    username: str
    password: str


@dataclass(frozen=True)
class LoginUserDTO:
    username: str
    password: str


@dataclass(frozen=True)
class AuthTokenDTO:
    access_token: str
    token_type: str = "Bearer"
