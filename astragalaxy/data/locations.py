import json

from pydantic import BaseModel


class Location(BaseModel):
    code: str
    emoji: str
    multiplayer: bool


def load_locations(file_path: str) -> list[Location]:
    with open(file_path) as f:
        raw_locations = json.load(f)

    return [Location.model_validate_json(loc) for loc in raw_locations]


LOCATIONS = load_locations("data/locations.json")
