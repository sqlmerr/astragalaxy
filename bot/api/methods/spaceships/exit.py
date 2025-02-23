from uuid import UUID

import api


class Exit:
    async def exit_my_spaceship(
        self: "api.Api", jwt_token: str, spaceship_id: UUID
    ) -> bool:
        return await self._change_sit_status(jwt_token, spaceship_id, "exit")
