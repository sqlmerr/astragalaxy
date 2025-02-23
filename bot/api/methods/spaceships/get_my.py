import api
from api.exceptions import APIError
from api.types.spaceship import Spaceship


class GetMy:
    async def get_my_spaceships(self: "api.Api", jwt_token: str) -> list[Spaceship]:
        response = await self.api.get(
            "/spaceships/my", headers={"Authorization": f"Bearer {jwt_token}"}, raw=True
        )

        json: list[dict] | dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            print(message)
            raise APIError(message=message, status_code=response.status_code)

        return [Spaceship.model_validate(sp) for sp in json]
