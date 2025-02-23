from uuid import UUID

import api
from api.exceptions import APIError
from api.types.system import System


class GetById:
    async def get_system_by_id(
        self: "api.Api", jwt_token: str, system_id: UUID | str
    ) -> System:
        response = await self.api.get(
            f"/systems/{system_id}",
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True,
        )
        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        return System.model_validate(json)
