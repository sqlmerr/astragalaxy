from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from .point import Point

class Station(Base):
    __tablename__ = "stations"
    id: Mapped[str] = mapped_column(primary_key=True)
    point: Mapped[Point] = relationship()
    point_id: Mapped[str] = mapped_column(ForeignKey("points.id"))

