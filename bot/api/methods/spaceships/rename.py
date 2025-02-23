from uuid import UUID

import api
from api.exceptions import APIError


class Rename:
    async def rename_my_spaceship(
        self: "api.Api", jwt_token: str, spaceship_id: UUID, name: str
    ) -> int:
        response = await self.api.put(
            "/spaceships/my/rename",
            headers={"Authorization": f"Bearer {jwt_token}"},
            json={"name": name, "spaceship_id": str(spaceship_id)},
            raw=True,
        )

        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        custom_status_code = json["custom_status_code"]
        return custom_status_code
