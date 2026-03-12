from dataclasses import dataclass

from astragalaxy.database.models.planet import PlanetThreatEnum, Planet


@dataclass(frozen=True)
class PlanetDTO:
    id: str
    name: str
    system_id: str
    threat: PlanetThreatEnum

    @classmethod
    def from_model(cls, model: Planet) -> "PlanetDTO":
        return cls(
            id=model.id, name=model.name, system_id=model.system_id, threat=model.threat
        )


@dataclass(frozen=True)
class CreatePlanetDTO:
    name: str
    system_id: str
    threat: PlanetThreatEnum
