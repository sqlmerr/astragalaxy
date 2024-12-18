from api.base import ApiBase
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

    async def register_user(self, user_id: int, username: str) -> User | None:
        response = await self.api.post(
            "/auth/register/telegram",
            json={"username": username, "telegram_id": user_id},
            headers={"secret-token": config.secret_token},
        )
        if isinstance(response, dict):
            return User.model_validate(response)

    # async def login_user
