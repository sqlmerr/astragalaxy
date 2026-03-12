from pydantic import BaseModel, Field


class PlanetNavigationSchema(BaseModel):
    planet_id: str


class HyperjumpSchema(BaseModel):
    path: str = Field(examples=["abcdef12->systemid->12345678"])
