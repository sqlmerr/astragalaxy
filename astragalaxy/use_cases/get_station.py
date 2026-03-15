from dataclasses import dataclass

from astragalaxy.dto.point import PointDTO
from astragalaxy.dto.station import StationDTO
from astragalaxy.exceptions.point import PointNotFound
from astragalaxy.exceptions.staion import StationNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.station.repo import StationRepo


@dataclass(frozen=True)
class GetStation:
    repo: StationRepo
    idp: IdentityProvider

    async def execute(self, station_id: str) -> StationDTO:
        await self.idp.get_current_user()

        station = await self.repo.find_one_station(station_id)
        if not station:
            raise StationNotFound

        return StationDTO.from_model(station)
