from uuid import UUID

from pydantic import BaseModel


class System(BaseModel):
    id: UUID
    name: str
