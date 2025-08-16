from enum import Enum

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from voidspace.database import Base
from .system import System


class PlanetThreatEnum(Enum):
    RADIATION = "radiation"
    TOXINS = "toxins"
    FREEZING = "freezing"
    HEAT = "heat"


class Planet(Base):
    __tablename__ = "planets"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(default="undefined")
    system: Mapped[System] = relationship()
    system_id: Mapped[str] = mapped_column(ForeignKey("systems.id"))
    threat: Mapped[PlanetThreatEnum]
