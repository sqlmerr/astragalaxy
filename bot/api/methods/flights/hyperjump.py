from uuid import UUID

import api


class HyperJump:
    async def hyperjump(
        self: "api.Api", jwt_token: str, system_id: UUID, spaceship_id: UUID
    ) -> int:
        response = await self.api.post(
            "/flights/hyperjump",
            json={"system_id": str(system_id), "spaceship_id": str(spaceship_id)},
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True
        )

        return response.status_code
