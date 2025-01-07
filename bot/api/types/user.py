from uuid import UUID

from pydantic import BaseModel

from api.types.spaceship import Spaceship


class User(BaseModel):
    id: UUID
    username: str
    telegram_id: int
    spaceships: list[Spaceship]
    in_spaceship: bool
    location_id: UUID
    system_id: UUID
