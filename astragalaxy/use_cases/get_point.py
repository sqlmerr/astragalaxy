from dataclasses import dataclass

from astragalaxy.dto.point import PointDTO
from astragalaxy.exceptions.point import PointNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.point.repo import PointRepo


@dataclass(frozen=True)
class GetPoint:
    repo: PointRepo
    idp: IdentityProvider

    async def execute(self, point_id: str) -> PointDTO:
        await self.idp.get_current_user()

        point = await self.repo.find_one_point(point_id)
        if not point:
            raise PointNotFound()

        return PointDTO.from_model(point)
