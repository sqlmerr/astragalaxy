from uuid import UUID

from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from voidspace.api.dependencies import JwtSecurity, CharacterSecurity
from voidspace.api.schemas import DataSchema
from voidspace.api.schemas.spaceship import SpaceshipSchema, RenameSpaceshipSchema
from voidspace.dto.spaceship import RenameSpaceshipDTO
from voidspace.use_cases.get_spaceship import GetCharacterSpaceships, GetSpaceship
from voidspace.use_cases.rename_spaceship import RenameSpaceship

router = APIRouter(prefix="/spaceships", route_class=DishkaRoute, tags=["Spaceships"])


@router.get("/my", dependencies=[JwtSecurity, CharacterSecurity])
async def get_my_spaceships(
    use_case: FromDishka[GetCharacterSpaceships],
) -> DataSchema[SpaceshipSchema]:
    spaceships = await use_case.execute()

    return DataSchema(data=[SpaceshipSchema.from_dto(sp) for sp in spaceships])


@router.post("/my/{id}/rename", dependencies=[JwtSecurity, CharacterSecurity])
async def rename_spaceship(
    id: UUID, data: RenameSpaceshipSchema, use_case: FromDishka[RenameSpaceship]
) -> SpaceshipSchema:
    spaceship = await use_case.execute(
        RenameSpaceshipDTO(name=data.name, spaceship_id=id)
    )

    return SpaceshipSchema.from_dto(spaceship)


@router.get("/{id}", dependencies=[JwtSecurity])
async def get_spaceship(
    id: UUID, use_case: FromDishka[GetSpaceship]
) -> SpaceshipSchema:
    """It can only be used to obtain a spaceship that belongs to one of the characters of an authorized user."""

    spaceship = await use_case.execute(id)
    return SpaceshipSchema.from_dto(spaceship)
