from dataclasses import dataclass
from uuid import UUID

from astragalaxy.dto.character import CharacterDTO
from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.exceptions.character import CharacterNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo


@dataclass(frozen=True)
class CharacterFilters:
    id: UUID | None
    code: str | None


@dataclass(frozen=True)
class GetCharacterFiltered:
    repo: CharacterRepo
    idp: IdentityProvider

    async def execute(self, data: CharacterFilters) -> CharacterDTO:
        if data.id is not None:
            character = await self.repo.find_one_character(data.id)
        elif data.code is not None:
            character = await self.repo.find_one_character_by_code(data.code)
        else:
            raise ValueError("Either id or code must be provided")

        if not character:
            raise CharacterNotFound()
        return CharacterDTO.from_model(character)


@dataclass(frozen=True)
class GetUserCharacters:
    repo: CharacterRepo
    idp: IdentityProvider

    async def execute(self) -> list[CharacterDTO]:
        current_user_id = self.idp.get_current_user_id()
        characters = await self.repo.find_all_characters_by_user_id(current_user_id)
        return [CharacterDTO.from_model(character) for character in characters]
