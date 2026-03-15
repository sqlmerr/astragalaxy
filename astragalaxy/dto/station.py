from dataclasses import dataclass

from astragalaxy.database.models import Station


@dataclass(frozen=True)
class StationDTO:
    id: str
    point_id: str

    @classmethod
    def from_model(cls, model: Station) -> "StationDTO":
        return cls(
            id=model.id, point_id=model.point_id
        )


@dataclass(frozen=True)
class CreateStationDTO:
    name: str
