from dataclasses import dataclass
from uuid import UUID

from jwt import PyJWTError
from starlette.datastructures import Headers

from astragalaxy.dto.character import CharacterDTO
from astragalaxy.dto.user import UserDTO
from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.exceptions.character import CharacterNotFound
from astragalaxy.exceptions.user import InvalidToken
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.identity_provider import IdentityProvider
from astragalaxy.interfaces.user.repo import UserRepo
from astragalaxy.jwt_token_processor import JwtTokenProcessor


@dataclass(frozen=True)
class IdentityProviderImpl(IdentityProvider):
    user_repo: UserRepo
    character_repo: CharacterRepo
    headers: Headers
    jwt_token_processor: JwtTokenProcessor

    def get_current_user_id(self) -> UUID:
        header = self.headers.get("Authorization")
        if not header:
            raise InvalidToken()

        token = header.split()[-1]
        try:
            payload = self.jwt_token_processor.decode(token)
        except PyJWTError:
            raise InvalidToken()

        try:
            user_id = UUID(payload["sub"])
        except ValueError:
            raise InvalidToken()

        return user_id

    async def get_current_user(self) -> UserDTO:
        user_id = self.get_current_user_id()

        user = await self.user_repo.find_one_user(user_id)
        if not user:
            raise InvalidToken()

        return UserDTO.from_model(user)

    def get_current_character_id(self) -> UUID:
        header = self.headers.get("X-Character-ID")
        if not header:
            raise CharacterNotFound()

        return UUID(header)

    async def get_current_character(self) -> CharacterDTO:
        character_id = self.get_current_character_id()

        character = await self.character_repo.find_one_character(character_id)
        if not character:
            raise CharacterNotFound()
        current_user_id = self.get_current_user_id()
        if character.user_id != current_user_id:
            raise AccessDeniedError()

        return CharacterDTO.from_model(character)
