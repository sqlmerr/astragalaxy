from pydantic import BaseModel

from astragalaxy.dto.station import StationDTO


class StationSchema(BaseModel):
    id: str
    point_id: str

    @classmethod
    def from_dto(cls, dto: StationDTO) -> "StationSchema":
        return cls.model_validate(dto, from_attributes=True)
