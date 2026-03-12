from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from .planet import Planet
from .system import System
from .user import User


class Character(Base):
    __tablename__ = "characters"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    code: Mapped[str] = mapped_column(unique=True)
    location: Mapped[str]
    in_spaceship: Mapped[bool] = mapped_column(default=False)
    user: Mapped[User] = relationship()
    user_id: Mapped[UUID] = mapped_column(ForeignKey("users.id"))
    system: Mapped[System] = relationship()
    system_id: Mapped[str] = mapped_column(ForeignKey("systems.id"))
    planet: Mapped[Planet] = relationship()
    planet_id: Mapped[str] = mapped_column(ForeignKey("planets.id"), nullable=True)
