from uuid import UUID

import api


class PlanetFlight:
    async def flight_to_planet(
        self: "api.Api", jwt_token: str, planet_id: UUID, spaceship_id: UUID
    ) -> bool:
        response = await self.api.post(
            "/flights/planet",
            json={"planet_id": str(planet_id), "spaceship_id": str(spaceship_id)},
            headers={"Authorization": f"Bearer {jwt_token}"},
            raw=True
        )

        if response.status_code != 200:
            return False

        json = response.json()
        return json["ok"]
