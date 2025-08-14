from dataclasses import dataclass
from uuid import UUID

from jwt import PyJWTError
from starlette.datastructures import Headers

from voidspace.dto.character import CharacterDTO
from voidspace.dto.user import UserDTO
from voidspace.exceptions.character import CharacterNotFound
from voidspace.exceptions.user import InvalidToken, UserNotFound
from voidspace.interfaces.character.read import CharacterReader
from voidspace.interfaces.user import UserReader
from voidspace.jwt_token_processor import JwtTokenProcessor


@dataclass(frozen=True)
class IdentityProvider:
    user_reader: UserReader
    character_reader: CharacterReader
    headers: Headers
    jwt_token_processor: JwtTokenProcessor

    def get_current_user_id(self) -> UUID:
        header = self.headers.get("Authorization")
        if not header:
            raise InvalidToken

        token = header.split()[-1]
        try:
            payload = self.jwt_token_processor.decode(token)
        except PyJWTError:
            raise InvalidToken

        try:
            user_id = UUID(payload["sub"])
        except ValueError:
            raise InvalidToken

        return user_id

    async def get_current_user(self) -> UserDTO:
        user_id = self.get_current_user_id()

        try:
            user = await self.user_reader.get_user_by_id(user_id)
        except UserNotFound:
            raise InvalidToken

        return user

    async def get_current_character(self) -> CharacterDTO:
        header = self.headers.get("X-Character-ID")
        if not header:
            raise CharacterNotFound

        character = await self.character_reader.get_character_by_id(UUID(header))
        return character
