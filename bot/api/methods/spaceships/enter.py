from uuid import UUID

import api


class Enter:
    async def enter_my_spaceship(
        self: "api.Api", jwt_token: str, spaceship_id: UUID
    ) -> bool:
        return await self._change_sit_status(jwt_token, spaceship_id, "enter")
