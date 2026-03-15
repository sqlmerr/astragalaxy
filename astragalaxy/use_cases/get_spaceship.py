from dataclasses import dataclass
from uuid import UUID

from astragalaxy.dto.spaceship import SpaceshipDTO
from astragalaxy.exceptions.spaceship import SpaceshipNotFoundError
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo


@dataclass(frozen=True)
class GetSpaceship:
    repo: SpaceshipRepo
    character_repo: CharacterRepo
    idp: IdentityProvider

    async def execute(self, data: UUID) -> SpaceshipDTO:
        current_user_id = self.idp.get_current_user_id()

        spaceship = await self.repo.find_one_by_id(data)
        if not spaceship:
            raise SpaceshipNotFoundError()

        spaceship_character = await self.character_repo.find_one_character(
            spaceship.character_id
        )
        if not spaceship_character:
            raise SpaceshipNotFoundError()

        if spaceship_character.user_id != current_user_id:
            raise SpaceshipNotFoundError()

        return SpaceshipDTO.from_model(spaceship)


@dataclass(frozen=True)
class GetActiveSpaceship:
    repo: SpaceshipRepo
    idp: IdentityProvider

    commiter: Commiter

    async def execute(self) -> SpaceshipDTO:
        character = await self.idp.get_current_character()

        spaceship = await self.repo.find_one_active_by_character(character.id)
        if spaceship:
            return SpaceshipDTO.from_model(spaceship)

        character_spaceships = await self.repo.find_all_by_character_id(character.id)
        if len(character_spaceships) == 0:
            raise SpaceshipNotFoundError()

        spaceship = character_spaceships[0]
        spaceship.active = True
        self.repo.add(spaceship)
        dto = SpaceshipDTO.from_model(spaceship)
        await self.commiter.commit()
        return dto


@dataclass(frozen=True)
class GetCharacterSpaceships:
    repo: SpaceshipRepo
    idp: IdentityProvider

    async def execute(self) -> list[SpaceshipDTO]:
        character = await self.idp.get_current_character()

        spaceships = await self.repo.find_all_by_character_id(character.id)
        return [SpaceshipDTO.from_model(s) for s in spaceships]
