from api.types.token import TokenPair
import api


class Login:
    async def login_user(
        self: "api.Api", user_id: int, token: str | None = None
    ) -> TokenPair | None:
        if not token:
            token = await self.get_user_token(user_id)
            if not token:
                await self.register_user(user_id)
                token = await self.get_user_token(user_id)

        jwt_token = await self.api.post("/auth/login", json={"token": token})

        return TokenPair(user_token=token, jwt_token=jwt_token["access_token"])
