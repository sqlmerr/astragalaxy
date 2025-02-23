import api
from api.exceptions import AuthError
from api.types.user import User


class GetMe:
    async def get_me(self: "api.Api", jwt_token: str) -> User:
        response = await self.api.get(
            "/auth/me", headers={"Authorization": f"Bearer {jwt_token}"}, raw=True
        )

        if response.status_code != 200:
            json: dict = response.json()
            message = json.get("message", None)
            print(message)
            raise AuthError(message=message, status_code=response.status_code)

        json: dict = response.json()
        user = User.model_validate(json)
        return user
