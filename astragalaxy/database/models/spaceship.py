from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from .character import Character
from .point import Point
from .system import System


class Spaceship(Base):
    __tablename__ = "spaceships"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    name: Mapped[str]
    character: Mapped[Character] = relationship()
    character_id: Mapped[UUID] = mapped_column(ForeignKey("characters.id"))
    active: Mapped[bool]
    point: Mapped[Point] = relationship()
    point_id: Mapped[str] = mapped_column(
        ForeignKey("points.id", ondelete="SET NULL", onupdate="CASCADE")
    )
