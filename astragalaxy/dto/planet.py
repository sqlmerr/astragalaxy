from dataclasses import dataclass

from astragalaxy.database.models.planet import PlanetThreatEnum, Planet


@dataclass(frozen=True)
class PlanetDTO:
    id: str
    name: str
    point_id: str
    threat: PlanetThreatEnum

    @classmethod
    def from_model(cls, model: Planet) -> "PlanetDTO":
        return cls(
            id=model.id, name=model.name, point_id=model.point_id, threat=model.threat
        )


@dataclass(frozen=True)
class CreatePlanetDTO:
    name: str
    point_id: str
    threat: PlanetThreatEnum
