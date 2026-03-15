from dataclasses import dataclass

from . import NotFoundError


@dataclass(eq=False)
class PointNotFound(NotFoundError):
    @property
    def message(self) -> str:
        return "Point not found"
