from dataclasses import dataclass
import random

from astragalaxy.cooldown_manager import CooldownManager
from astragalaxy.database.models import System
from astragalaxy.dto.cooldown import SetCooldownDTO, CooldownDTO
from astragalaxy.exceptions import AppError, CharacterInCooldown
from astragalaxy.exceptions.spaceship import (
    CharacterNeedsToBeInSpaceship,
    InvalidHyperjumpPath,
)
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.character.repo import CharacterRepo
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.spaceship.repo import SpaceshipRepo
from astragalaxy.interfaces.station.repo import StationRepo
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.interfaces.system_connection.repo import SystemConnectionRepo


def parse_path(path: str) -> list[str]:
    return path.split("->")


@dataclass(frozen=True)
class Hyperjump:
    spaceship_repo: SpaceshipRepo
    connection_repo: SystemConnectionRepo
    system_repo: SystemRepo
    character_repo: CharacterRepo
    point_repo: PointRepo
    station_repo: StationRepo
    idp: IdentityProvider
    commiter: Commiter

    cooldown_manager: CooldownManager

    async def _process_path(
        self, current_system_id: str, path: list[str]
    ) -> list[System]:
        # [B, C, D]
        system_1 = await self.system_repo.find_one_system(path[0])  # B
        if not system_1:
            raise InvalidHyperjumpPath()
        conns_1 = await self.connection_repo.find_connections_by_system_ids(
            current_system_id, system_1.id
        )  # A -> B
        if len(conns_1) == 0:
            raise InvalidHyperjumpPath()

        latest_system_id = system_1.id
        systems = [system_1]
        for system_id in path[1:]:  # [C, D]
            system = await self.system_repo.find_one_system(system_id)  # B, C
            if not system:
                raise InvalidHyperjumpPath()
            conns = await self.connection_repo.find_connections_by_system_ids(
                system_id, latest_system_id
            )  # B -> C, C -> D
            if len(conns) == 0:
                raise InvalidHyperjumpPath
            systems.append(system)
            latest_system_id = system.id

        return systems

    async def execute(self, path: str) -> CooldownDTO:
        parsed_path = parse_path(path)
        current_character_id = self.idp.get_current_character_id()
        current_character = await self.character_repo.find_one_character(
            current_character_id
        )
        if not current_character:
            raise AppError()

        spaceship = await self.spaceship_repo.find_one_active_by_character(
            current_character_id
        )
        if not spaceship:
            raise AppError()

        cooldown = await self.cooldown_manager.get(current_character_id)
        if cooldown.remaining_seconds > 0:
            raise CharacterInCooldown()

        if not current_character.in_spaceship:
            raise CharacterNeedsToBeInSpaceship()

        point = await self.point_repo.find_one_point(current_character.point_id)
        if not point:
            raise AppError()

        systems = await self._process_path(point.system_id, parsed_path)
        destination_system = systems[-1]
        seconds = len(systems) * 60

        stations = await self.station_repo.get_stations_by_system(destination_system.id)
        if len(stations) == 0:
            raise AppError()
        station = random.choice(stations)

        cooldown = await self.cooldown_manager.set(
            SetCooldownDTO(
                character_id=current_character_id,
                seconds=seconds,
                action="navigation_hyperjump",
            )
        )
        
        current_character.point_id = station.point_id
        spaceship.point_id = station.point_id

        self.character_repo.add(current_character)
        self.spaceship_repo.add(spaceship)
        await self.commiter.commit()

        return cooldown
