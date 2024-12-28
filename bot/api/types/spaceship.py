from pydantic import BaseModel


class Spaceship(BaseModel):
    id: str
    name: str
    user_id: str
    location_id: str
