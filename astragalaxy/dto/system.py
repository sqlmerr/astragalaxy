from dataclasses import dataclass


@dataclass(frozen=True)
class CreateSystemDTO:
    name: str
    connections: list[str]


@dataclass(frozen=True)
class SystemDTO:
    id: str
    name: str
    connections: list[str]
