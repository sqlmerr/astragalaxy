from pydantic import BaseModel

from voidspace.dto.cooldown import CooldownDTO


class CooldownSchema(BaseModel):
    remaining_seconds: int
    total_seconds: int
    action: str

    @classmethod
    def from_dto(cls, dto: CooldownDTO) -> "CooldownSchema":
        return cls(
            remaining_seconds=int(dto.remaining_seconds),
            total_seconds=int(dto.seconds),
            action=dto.action,
        )
