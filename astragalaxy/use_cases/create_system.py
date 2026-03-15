from dataclasses import dataclass
from uuid import uuid4

from astragalaxy.database.models import Point, Station, System
from astragalaxy.database.models.system_connection import SystemConnection
from astragalaxy.dto.system import CreateSystemDTO, SystemDTO
from astragalaxy.exceptions import AccessDeniedError, AppError
from astragalaxy.exceptions.system import SystemNotFound
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.point.repo import PointRepo
from astragalaxy.interfaces.session import Commiter
from astragalaxy.interfaces.station.repo import StationRepo
from astragalaxy.interfaces.system.repo import SystemRepo
from astragalaxy.interfaces.system_connection.repo import SystemConnectionRepo
from astragalaxy.utils import generate_random_id


@dataclass(frozen=True)
class CreateSystem:
    repo: SystemRepo
    connection_repo: SystemConnectionRepo
    point_repo: PointRepo
    station_repo: StationRepo
    idp: IdentityProvider
    commiter: Commiter

    async def execute(self, data: CreateSystemDTO) -> SystemDTO:
        current_user = await self.idp.get_current_user()
        if current_user.username != "admin":  # TODO: roles
            raise AccessDeniedError()

        system_id = generate_random_id(8)
        self.repo.add(System(id=system_id, name=data.name))

        system = await self.repo.find_one_system(system_id)
        if not system:
            raise AppError()
        for conn in data.connections:
            sys = await self.repo.find_one_system(conn)
            if not sys:
                raise SystemNotFound()
            self.connection_repo.add_connection(
                SystemConnection(
                    id=uuid4(), system_from_id=system.id, system_to_id=sys.id
                )
            )
            self.connection_repo.add_connection(
                SystemConnection(
                    id=uuid4(), system_from_id=sys.id, system_to_id=system.id
                )
            )

        station_point = Point(id=generate_random_id(16), name="station_point", system_id=system_id)
        self.point_repo.add(station_point)
        station = Station(id=generate_random_id(8), point_id=station_point.id)
        self.station_repo.add(station)
        
        dto = SystemDTO(
            id=system.id,
            connections=data.connections,
            name=system.name,
        )
        
        await self.commiter.commit()

        return dto
