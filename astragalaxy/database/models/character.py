from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from .point import Point
from .system import System
from .user import User


class Character(Base):
    __tablename__ = "characters"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    code: Mapped[str] = mapped_column(unique=True)
    in_spaceship: Mapped[bool] = mapped_column(default=False)
    user: Mapped[User] = relationship()
    user_id: Mapped[UUID] = mapped_column(ForeignKey("users.id"))
    point: Mapped[Point] = relationship()
    point_id: Mapped[str] = mapped_column(ForeignKey("points.id"))
