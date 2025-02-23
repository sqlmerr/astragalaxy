import api
from config_reader import config


class GetToken:
    async def get_user_token(self: "api.Api", user_id: int) -> str | None:
        token_res = await self.api.get(
            "/auth/token/sudo",
            params={"telegram_id": user_id},
            headers={"secret-token": config.secret_token},
            raw=True,
        )

        if token_res.status_code != 200:
            return None

        token: str = str(user_id) + ":" + (token_res.json())["token"]

        return token
