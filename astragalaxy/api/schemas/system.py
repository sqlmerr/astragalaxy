from pydantic import BaseModel

from astragalaxy.dto.system import SystemDTO


class SystemSchema(BaseModel):
    id: str
    name: str
    connections: list[str]

    @classmethod
    def from_dto(cls, dto: SystemDTO) -> "SystemSchema":
        return cls.model_validate(dto, from_attributes=True)
