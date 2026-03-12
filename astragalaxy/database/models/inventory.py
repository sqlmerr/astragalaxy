import enum
from uuid import UUID, uuid4

from sqlalchemy.orm import Mapped, mapped_column

from astragalaxy.database import Base


class InventoryOwnerEnum(enum.Enum):
    CHARACTER = "character"
    SPACESHIP = "spaceship"


class Inventory(Base):
    __tablename__ = "inventories"
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    owner: Mapped[InventoryOwnerEnum]
    owner_id: Mapped[UUID]
