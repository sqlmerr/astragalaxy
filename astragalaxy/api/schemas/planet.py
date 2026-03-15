from pydantic import BaseModel

from astragalaxy.database.models.planet import PlanetThreatEnum
from astragalaxy.dto.planet import PlanetDTO, CreatePlanetDTO


class PlanetSchema(BaseModel):
    id: str
    name: str
    point_id: str
    threat: PlanetThreatEnum

    @classmethod
    def from_dto(cls, dto: PlanetDTO) -> "PlanetSchema":
        return cls.model_validate(dto, from_attributes=True)


class CreatePlanetSchema(BaseModel):
    name: str
    point_id: str
    threat: PlanetThreatEnum

    def into_dto(self) -> "CreatePlanetDTO":
        return CreatePlanetDTO(name=self.name, point_id=self.point_id, threat=self.threat)
