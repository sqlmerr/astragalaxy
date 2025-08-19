from uuid import UUID, uuid4

from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column, relationship

from voidspace.database import Base
from voidspace.database.models.inventory import Inventory


class Resource(Base):
    __tablename__ = "resources"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    code: Mapped[str]
    quantity: Mapped[int]
    inventory: Mapped[Inventory] = relationship()
    inventory_id: Mapped[UUID] = mapped_column(ForeignKey("inventories.id"))
