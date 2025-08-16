from dataclasses import dataclass
from uuid import UUID

from voidspace.dto.character import CharacterDTO
from voidspace.exceptions import AccessDeniedError
from voidspace.exceptions.character import CharacterNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class CharacterFilters:
    id: UUID | None
    code: str | None


@dataclass(frozen=True)
class GetCharacterFiltered(BaseUseCase[CharacterFilters, CharacterDTO]):
    repo: CharacterRepo
    idp: IdentityProvider

    async def execute(self, data: CharacterFilters) -> CharacterDTO:
        if data.id is not None:
            character = await self.repo.find_one_character(data.id)
        else:
            character = await self.repo.find_one_character_by_code(data.code)

        if not character:
            raise CharacterNotFound()

        current_user_id = self.idp.get_current_user_id()
        if character.user_id != current_user_id:
            raise AccessDeniedError()

        return CharacterDTO.from_model(character)


@dataclass(frozen=True)
class GetUserCharacters:
    repo: CharacterRepo
    idp: IdentityProvider

    async def execute(self) -> list[CharacterDTO]:
        current_user_id = self.idp.get_current_user_id()
        characters = await self.repo.find_all_characters_by_user_id(current_user_id)
        return [CharacterDTO.from_model(character) for character in characters]
