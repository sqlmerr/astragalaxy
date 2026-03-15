from typing import Protocol


class Commiter(Protocol):
    async def commit(self) -> None:
        raise NotImplementedError


class Refresher(Protocol):
    async def refresh(self, instance: object) -> None:
        raise NotImplementedError