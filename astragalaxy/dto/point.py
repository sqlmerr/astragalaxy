from dataclasses import dataclass

from astragalaxy.database.models import Point


@dataclass(frozen=True)
class PointDTO:
    id: str
    name: str
    system_id: str

    @classmethod
    def from_model(cls, model: Point) -> "PointDTO":
        return cls(
            id=model.id,
            name=model.name,
            system_id=model.system_id
        )


@dataclass(frozen=True)
class CreatePointDTO:
    name: str
    system_id: str
