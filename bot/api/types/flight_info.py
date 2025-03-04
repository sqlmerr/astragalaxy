import uuid

from pydantic import BaseModel


class FlyInfo(BaseModel):
    flying: bool
    destination: str
    destination_id: uuid.UUID
    remaining_time: int
    flown_out_at: int
