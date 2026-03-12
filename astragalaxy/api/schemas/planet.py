from pydantic import BaseModel

from astragalaxy.database.models.planet import PlanetThreatEnum
from astragalaxy.dto.planet import PlanetDTO


class PlanetSchema(BaseModel):
    id: str
    name: str
    system_id: str
    threat: PlanetThreatEnum

    @classmethod
    def from_dto(cls, dto: PlanetDTO) -> "PlanetSchema":
        return cls(id=dto.id, name=dto.name, system_id=dto.system_id, threat=dto.threat)
