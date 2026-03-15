from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from .system import System


class Point(Base):
    __tablename__ = "points"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(default="undefined")
    system: Mapped[System] = relationship()
    system_id: Mapped[str] = mapped_column(ForeignKey("systems.id"))
