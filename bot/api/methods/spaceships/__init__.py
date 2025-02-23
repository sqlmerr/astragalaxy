import api

from uuid import UUID

from . import get_my, exit, enter, rename
from api.exceptions import APIError


class Spaceships(get_my.GetMy, exit.Exit, enter.Enter, rename.Rename):
    async def _change_sit_status(
        self: "api.Api", jwt_token: str, spaceship_id: UUID, status: str
    ) -> bool:
        response = await self.api.post(
            f"/spaceships/my/{spaceship_id}/{status}",
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True,
        )

        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        ok = json["ok"]
        return ok
