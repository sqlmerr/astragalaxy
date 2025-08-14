from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from ..dependencies import JwtSecurity, CharacterSecurity
from ..schemas import DataSchema
from ..schemas.character import CreateCharacterSchema, CharacterSchema
from ...dto.character import CreateCharacterDTO
from ...identity_provider import IdentityProvider
from ...interfaces.character.create import CharacterWriter
from ...interfaces.character.read import CharacterReader

router = APIRouter(prefix="/characters", route_class=DishkaRoute, tags=["Characters"])


@router.post("/", dependencies=[JwtSecurity], status_code=201)
async def create_character(
    data: CreateCharacterSchema,
    character_writer: FromDishka[CharacterWriter],
    identity_provider: FromDishka[IdentityProvider],
) -> CharacterSchema:
    user_id = identity_provider.get_current_user_id()
    character = await character_writer.create_character(
        CreateCharacterDTO(user_id=user_id, code=data.code)
    )

    return CharacterSchema.from_dto(character)v


@router.get("/my", dependencies=[JwtSecurity])
async def get_my_characters(
    character_reader: FromDishka[CharacterReader],
    identity_provider: FromDishka[IdentityProvider],
) -> DataSchema[CharacterSchema]:
    current_user = await identity_provider.get_current_user()
    characters = await character_reader.get_characters_by_user(current_user.id)

    schemas = [CharacterSchema.from_dto(c) for c in characters]

    return DataSchema(data=schemas)


@router.get("/", dependencies=[JwtSecurity, CharacterSecurity])
async def get_current_character(
    identity_provider: FromDishka[IdentityProvider],
) -> CharacterSchema:
    character = await identity_provider.get_current_character()
    return CharacterSchema.from_dto(character)
