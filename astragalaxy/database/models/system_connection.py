from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from .system import System
from .. import Base


class SystemConnection(Base):
    __tablename__ = "system_connections"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    system_from_id: Mapped[str] = mapped_column(ForeignKey("systems.id"))
    system_to_id: Mapped[str] = mapped_column(ForeignKey("systems.id"))

    system_from: Mapped[System] = relationship(foreign_keys=[system_from_id])
    system_to: Mapped[System] = relationship(foreign_keys=[system_to_id])
