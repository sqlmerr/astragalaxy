from dataclasses import dataclass

from astragalaxy.dto.station import StationDTO
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.station.repo import StationRepo


@dataclass(frozen=True)
class GetPointStations:
    repo: StationRepo
    idp: IdentityProvider

    async def execute(self, point_id: str) -> list[StationDTO]:
        await self.idp.get_current_user()

        stations = await self.repo.get_stations_by_point(point_id)

        return [StationDTO.from_model(s) for s in stations]
