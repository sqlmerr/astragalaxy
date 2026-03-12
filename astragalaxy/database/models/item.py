from typing import Any
from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from astragalaxy.database import Base
from astragalaxy.database.models.inventory import Inventory


class Item(Base):
    __tablename__ = "items"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    code: Mapped[str]
    inventory_id: Mapped[UUID] = mapped_column(ForeignKey("inventories.id"))
    inventory: Mapped[Inventory] = relationship()
    data: Mapped[dict[str, Any]]
    durability: Mapped[int] = mapped_column(nullable=True, default=100)
