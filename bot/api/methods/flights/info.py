from uuid import UUID

import api
from api.exceptions import APIError
from api.types.flight_info import FlyInfo


class GetFlightInfo:
    async def get_flight_info(
        self: "api.Api", jwt_token: str, spaceship_id: UUID
    ) -> FlyInfo:
        response = await self.api.get(
            "/flights/info?id={}".format(spaceship_id),
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True
        )
        if response.status_code != 200:
            raise APIError(response.json().get("message"), status_code=response.status_code)

        return FlyInfo.model_validate(response.json())
