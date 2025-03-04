from enum import StrEnum
from uuid import UUID

from pydantic import BaseModel


class PlanetThreat(StrEnum):
    RADIATION = "RADIATION"
    TOXINS = "TOXINS"
    FREEZING = "FREEZING"
    HEAT = "HEAT"


class Planet(BaseModel):
    id: UUID
    name: str
    system_id: UUID
    threat: PlanetThreat
