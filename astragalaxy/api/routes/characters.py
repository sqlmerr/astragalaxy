from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from ..dependencies import JwtSecurity, CharacterSecurity
from ..schemas import DataSchema
from ..schemas.character import CreateCharacterSchema, CharacterSchema
from ..schemas.cooldown import CooldownSchema
from ...cooldown_manager import CooldownManager
from ...dto.character import CreateCharacterDTO
from ...identity_provider import IdentityProvider
from ...use_cases.create_character import CreateCharacter
from ...use_cases.get_character import GetUserCharacters

router = APIRouter(prefix="/characters", route_class=DishkaRoute, tags=["Characters"])


@router.post("/", dependencies=[JwtSecurity], status_code=201)
async def create_character(
    data: CreateCharacterSchema,
    use_case: FromDishka[CreateCharacter],
) -> CharacterSchema:
    character = await use_case.execute(CreateCharacterDTO(code=data.code))

    return CharacterSchema.from_dto(character)


@router.get("/my", dependencies=[JwtSecurity])
async def get_my_characters(
    use_case: FromDishka[GetUserCharacters],
) -> DataSchema[CharacterSchema]:
    characters = await use_case.execute()
    schemas = [CharacterSchema.from_dto(c) for c in characters]

    return DataSchema(data=schemas)


@router.get("/current", dependencies=[JwtSecurity, CharacterSecurity])
async def get_current_character(
    identity_provider: FromDishka[IdentityProvider],
) -> CharacterSchema:
    character = await identity_provider.get_current_character()
    return CharacterSchema.from_dto(character)


@router.get("/current/cooldown", dependencies=[JwtSecurity, CharacterSecurity])
async def get_cooldown(
    cooldown_manager: FromDishka[CooldownManager],
    identity_provider: FromDishka[IdentityProvider],
) -> CooldownSchema:
    current_character = await identity_provider.get_current_character()
    cooldown = await cooldown_manager.get(current_character.id)

    return CooldownSchema.from_dto(cooldown)
