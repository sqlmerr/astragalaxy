from uuid import UUID

import api


class PlanetFlight:
    async def flight_to_planet(
        self: "api.Api", jwt_token: str, planet_id: UUID, spaceship_id: UUID
    ) -> int:
        response = await self.api.post(
            "/flights/planet",
            json={"planet_id": str(planet_id), "spaceship_id": str(spaceship_id)},
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True
        )

        return response.status_code
