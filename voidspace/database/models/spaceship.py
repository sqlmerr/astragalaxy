from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from .character import Character
from .system import System
from .planet import Planet
from voidspace.database import Base


class Spaceship(Base):
    __tablename__ = "spaceships"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    name: Mapped[str]
    location: Mapped[str]
    character: Mapped[Character] = relationship()
    character_id: Mapped[UUID] = mapped_column(ForeignKey("characters.id"))
    system: Mapped[System] = relationship()
    system_id: Mapped[UUID] = mapped_column(ForeignKey("systems.id"))
    planet: Mapped[Planet] = relationship()
    planet_id: Mapped[UUID] = mapped_column(ForeignKey("planets.id", ondelete="SET NULL", onupdate="CASCADE"), nullable=True)
