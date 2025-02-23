import api
from api.types.user import User
from config_reader import config


class Register:
    async def register_user(self: "api.api.Api", user_id: int) -> User | None:
        response = await self.api.post(
            "/auth/register",
            json={"username": f"id{user_id}", "telegram_id": user_id},
            headers={"secret-token": config.secret_token},
        )
        if isinstance(response, dict):
            return User.model_validate(response)
