from enum import Enum

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from .point import Point


class PlanetThreatEnum(Enum):
    RADIATION = "radiation"
    TOXINS = "toxins"
    FREEZING = "freezing"
    HEAT = "heat"


class Planet(Base):
    __tablename__ = "planets"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(default="undefined")
    point: Mapped[Point] = relationship()
    point_id: Mapped[str] = mapped_column(ForeignKey("points.id"))
    threat: Mapped[PlanetThreatEnum]
