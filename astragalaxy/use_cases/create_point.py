from dataclasses import dataclass

from astragalaxy.database.models import Point
from astragalaxy.dto.point import CreatePointDTO, PointDTO
from astragalaxy.exceptions import AccessDeniedError, AppError
from astragalaxy.exceptions.point import PointNotFound
from astragalaxy.exceptions.system import SystemNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.utils import generate_random_id


@dataclass(frozen=True)
class CreatePoint:
    repo: PointRepo
    system_repo: SystemRepo
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, data: CreatePointDTO) -> PointDTO:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":  # TODO: roles
            raise AccessDeniedError()

        sys = await self.system_repo.find_one_system(data.system_id)
        if not sys:
            raise SystemNotFound()

        point = Point(
            id=generate_random_id(16),
            name=data.name,
            system_id=data.system_id
        )
        self.repo.add(point)
        dto = PointDTO.from_model(point)
        await self.commiter.commit()

        return dto
