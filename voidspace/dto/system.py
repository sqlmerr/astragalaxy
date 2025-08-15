from dataclasses import dataclass

from voidspace.database.models import System


@dataclass(frozen=True)
class CreateSystemDTO:
    name: str
    locations: list[str]


@dataclass(frozen=True)
class SystemDTO:
    id: str
    name: str
    locations: list[str]

    @classmethod
    def from_model(cls, model: System) -> "SystemDTO":
        return cls(id=model.id, name=model.name, locations=model.locations)
