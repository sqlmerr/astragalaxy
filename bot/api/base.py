import aiohttp

from typing import Any
from config_reader import config


class ApiBase:
    def __init__(self) -> None:
        self.url = config.api_url

    async def request(self, method: str, path: str, raw: bool, **kwargs) -> Any | dict:
        async with aiohttp.ClientSession(self.url) as session:
            async with session.request(method, path, **kwargs) as response:
                if raw:
                    return response
                try:
                    return await response.json()
                except aiohttp.ContentTypeError:
                    return response.content

    async def post(
        self,
        path: str,
        json: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Any | dict:
        return await self.request(
            "POST", path, json=json, raw=raw, **kwargs
        )

    async def get(
        self,
        path: str,
        params: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Any | dict:
        return await self.request(
            "GET", path, params=params, raw=raw, **kwargs
        )

    async def put(
        self,
        path: str,
        data: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Any | dict:
        return await self.request(
            "PUT", path, data=data, raw=raw, **kwargs
        )

    async def delete(
        self,
        path: str,
        params: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Any | dict:
        return await self.request(
            "DELETE", path, params=params, raw=raw, **kwargs
        )