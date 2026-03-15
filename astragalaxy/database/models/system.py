from sqlalchemy import String
from sqlalchemy.dialects.postgresql import ARRAY
from sqlalchemy.orm import Mapped, mapped_column

from astragalaxy.database import Base


class System(Base):
    __tablename__ = "systems"
    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str]

