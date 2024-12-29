from api.base import ApiBase
from api.exceptions import AuthError, APIError
from api.types.spaceship import Spaceship
from api.types.token import TokenPair
from api.types.user import User
from config_reader import config


class Api:
    def __init__(self, api: ApiBase):
        self.api = api

    async def ping(self) -> bool:
        response = await self.api.get("/")
        if isinstance(response, dict) & response.get("ok", False):
            return True
        return False

    async def register_user(self, user_id: int) -> User | None:
        response = await self.api.post(
            "/auth/register/telegram",
            json={"username": f"id{user_id}", "telegram_id": user_id},
            headers={"secret-token": config.secret_token},
        )
        if isinstance(response, dict):
            return User.model_validate(response)

    async def get_user_token(self, user_id: int) -> str | None:
        token_res = await self.api.get(
            "/auth/token/sudo",
            params={"telegram_id": user_id},
            headers={"secret-token": config.secret_token},
            raw=True,
        )

        if token_res.status_code != 200:
            return

        token: str = str(user_id) + ":" + (token_res.json())["token"]

        return token

    async def login_user(
        self, user_id: int, token: str | None = None
    ) -> TokenPair | None:
        if not token:
            token = await self.get_user_token(user_id)
            if not token:
                await self.register_user(user_id)
                token = await self.get_user_token(user_id)

        jwt_token = await self.api.post("/auth/login", json={"token": token})

        return TokenPair(user_token=token, jwt_token=jwt_token["access_token"])

    async def get_me(self, jwt_token: str) -> User:
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

    async def get_my_spaceship(self, jwt_token: str) -> Spaceship:
        response = await self.api.get(
            "/spaceships/my", headers={"Authorization": f"Bearer {jwt_token}"}, raw=True
        )

        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            print(message)
            raise APIError(message=message, status_code=response.status_code)

        spaceship = Spaceship.model_validate(json)
        return spaceship

    async def get_out_of_my_spaceship(self, jwt_token: str) -> bool:
        response = await self.api.post(
            "/spaceships/my/getOut",
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True,
        )

        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        ok = json["ok"]
        return ok

    async def enter_my_spaceship(self, jwt_token: str) -> bool:
        response = await self.api.post(
            "/spaceships/my/enter",
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True,
        )

        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        ok = json["ok"]
        return ok

    async def rename_my_spaceship(self, jwt_token: str, name: str) -> int:
        response = await self.api.post("/spaceships/my/rename", headers={"Authorization": f"Bearer {jwt_token}"}, json={"name": name}, raw=True)

        json: dict = response.json()
        if response.status_code != 200:
            message = json.get("message", None)
            raise APIError(message=message, status_code=response.status_code)

        custom_status_code = json["custom_status_code"]
        return custom_status_code
