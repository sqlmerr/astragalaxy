from dataclasses import dataclass

from . import NotFoundError


@dataclass(eq=False)
class PlanetNotFound(NotFoundError):
    @property
    def message(self) -> str:
        return "Planet not found"
