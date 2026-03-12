from typing import Protocol
from uuid import UUID

from astragalaxy.dto.character import CharacterDTO
from astragalaxy.dto.user import UserDTO


class IdentityProvider(Protocol):
    def get_current_user_id(self) -> UUID:
        raise NotImplementedError

    async def get_current_user(self) -> UserDTO:
        raise NotImplementedError

    def get_current_character_id(self) -> UUID:
        raise NotImplementedError

    async def get_current_character(self) -> CharacterDTO:
        raise NotImplementedError
