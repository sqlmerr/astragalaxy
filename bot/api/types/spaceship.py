from uuid import UUID

from pydantic import BaseModel


class Spaceship(BaseModel):
    id: UUID
    name: str
    user_id: UUID
    location: str
    system_id: UUID
    planet_id: UUID
    player_sit_in: bool
