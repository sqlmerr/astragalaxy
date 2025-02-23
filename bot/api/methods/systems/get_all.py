import api
from api.exceptions import APIError

from api.types.system import System


class GetAll:
    async def get_all_systems(self: "api.Api", jwt_token: str) -> list[System]:
        response = await self.api.get(
            "/systems",
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True,
        )
        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        return [System.model_validate(sp) for sp in json]
