from dataclasses import dataclass

from . import NotFoundError


@dataclass(eq=False)
class SystemNotFound(NotFoundError):
    @property
    def message(self) -> str:
        return "System not found"
