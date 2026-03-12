from dataclasses import dataclass
from uuid import UUID

from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.exceptions.character import CharacterNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo


@dataclass(frozen=True)
class DeleteCharacter:
    repo: CharacterRepo
    idp: IdentityProvider

    async def execute(self, id: UUID) -> None:
        current_user_id = self.idp.get_current_user_id()
        character = await self.repo.find_one_character(id)
        if not character:
            raise CharacterNotFound()

        if character.user_id != current_user_id:
            raise AccessDeniedError()

        await self.repo.delete_character(id)
