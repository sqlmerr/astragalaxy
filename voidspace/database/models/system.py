from uuid import UUID, uuid4

from sqlalchemy import String
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.dialects.postgresql import ARRAY

from voidspace.database import Base


class System(Base):
    __tablename__ = "systems"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str]
    locations: Mapped[list[str]] = mapped_column(ARRAY(String))
