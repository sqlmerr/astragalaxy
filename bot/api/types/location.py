from uuid import UUID

from pydantic import BaseModel


class Location(BaseModel):
    id: UUID
    code: str
    multiplayer: bool
