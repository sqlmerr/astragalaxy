from uuid import UUID

import api
from api.exceptions import APIError
from api.types.planet import Planet


class GetPlanets:
    async def get_system_planets(
        self: "api.Api", jwt_token: str, system_id: UUID | str
    ) -> list[Planet]:
        response = await self.api.get(
            f"/systems/{system_id}/planets",
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True,
        )
        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        return [Planet.model_validate(sp) for sp in json]
