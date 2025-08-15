from typing import Protocol


class SystemDeleter(Protocol):
    async def delete_system(self, id: str) -> None:
        raise NotImplementedError
