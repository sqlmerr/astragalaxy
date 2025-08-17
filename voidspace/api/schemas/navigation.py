from pydantic import BaseModel


class PlanetNavigationSchema(BaseModel):
    planet_id: str
