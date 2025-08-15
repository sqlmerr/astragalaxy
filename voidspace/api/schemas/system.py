from pydantic import BaseModel

from voidspace.dto.system import SystemDTO


class SystemSchema(BaseModel):
    id: str
    name: str
    locations: list[str]

    @classmethod
    def from_dto(cls, dto: SystemDTO) -> "SystemSchema":
        return cls(id=dto.id, name=dto.name, locations=dto.locations)
