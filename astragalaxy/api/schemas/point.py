from pydantic import BaseModel

from astragalaxy.dto.point import PointDTO, CreatePointDTO


class CreatePointSchema(BaseModel):
    name: str
    system_id: str

    def into_dto(self) -> "CreatePointDTO":
        return CreatePointDTO(name=self.name, system_id=self.system_id)


class PointSchema(BaseModel):
    id: str
    name: str
    system_id: str

    @classmethod
    def from_dto(cls, dto: PointDTO) -> "PointSchema":
        return cls.model_validate(dto, from_attributes=True)
