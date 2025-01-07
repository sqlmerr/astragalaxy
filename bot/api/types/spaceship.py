from uuid import UUID

from pydantic import BaseModel


class Spaceship(BaseModel):
    id: UUID
    name: str
    user_id: UUID
    location_id: UUID
    flown_out_at: int
    flying: bool
    system_id: UUID
    planet_id: UUID
    player_sit_in: bool
