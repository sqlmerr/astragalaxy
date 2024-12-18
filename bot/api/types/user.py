from pydantic import BaseModel


class User(BaseModel):
    id: str
    username: str
    telegram_id: int
    spaceship_id: str | None
    in_spaceship: bool
    location_id: str
    system_id: str
