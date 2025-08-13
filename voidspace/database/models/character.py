from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from .user import User
from voidspace.database import Base


class Character(Base):
    __tablename__ = "characters"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    code: Mapped[str] = mapped_column(unique=True)
    location: Mapped[str]
    user: Mapped[User] = relationship()
    user_id: Mapped[UUID] = mapped_column(ForeignKey("users.id"))
