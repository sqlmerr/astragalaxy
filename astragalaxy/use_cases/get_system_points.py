from dataclasses import dataclass

from astragalaxy.dto.point import PointDTO
from astragalaxy.exceptions.point import PointNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.point.repo import PointRepo


@dataclass(frozen=True)
class GetSystemPoints:
    repo: PointRepo
    idp: IdentityProvider

    async def execute(self, system_id: str) -> list[PointDTO]:
        await self.idp.get_current_user()

        points = await self.repo.find_all_points_by_system(system_id=system_id)

        return [PointDTO.from_model(p) for p in points]
