from api.base import ApiBase
from .methods import Methods


class Api(Methods):
    def __init__(self, api: ApiBase):
        self.api = api

    async def ping(self) -> bool:
        response = await self.api.get("/")
        if isinstance(response, dict) & response.get("ok", False):
            return True
        return False
