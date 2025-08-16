from dataclasses import dataclass
from uuid import UUID

from voidspace.exceptions import AccessDeniedError
from voidspace.exceptions.character import CharacterNotFound
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.character.repo import CharacterRepo
from voidspace.use_cases import BaseUseCase


@dataclass(frozen=True)
class DeleteCharacter(BaseUseCase[UUID, None]):
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
