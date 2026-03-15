from dataclasses import dataclass

from . import NotFoundError


@dataclass(eq=False)
class StationNotFound(NotFoundError):
    @property
    def message(self) -> str:
        return "Station not found"
