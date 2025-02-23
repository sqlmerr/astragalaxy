from enum import StrEnum, auto
from uuid import UUID

from pydantic import BaseModel


class PlanetThreat(StrEnum):
    RADIATION = auto()
    TOXINS = auto()
    FREEZING = auto()
    HEAT = auto()


class Planet(BaseModel):
    id: UUID
    name: str
    system_id: UUID
    threat: PlanetThreat
