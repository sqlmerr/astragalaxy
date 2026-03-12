from dataclasses import dataclass

from . import NotFoundError


@dataclass(eq=False)
class InventoryNotFoundError(NotFoundError):
    @property
    def message(self) -> str:
        return "Inventory not found"
