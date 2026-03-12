import enum
import json

from pydantic import BaseModel


class ItemRarity(enum.Enum):
    COMMON = "common"
    RARE = "rare"
    LEGENDARY = "legendary"
    IMMORTAL = "immortal"


class Item(BaseModel):
    code: str
    emoji: str
    damage_per_use: int
    rarity: ItemRarity
    action: str


def load_items(file_path: str) -> list[Item]:
    with open(file_path) as f:
        raw = json.load(f)

    return [Item.model_validate_json(i) for i in raw]


ITEMS = load_items(file_path="data/items.json")
