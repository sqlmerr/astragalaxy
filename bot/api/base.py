from typing import Any

import httpx
from httpx import Response

from config_reader import config


class ApiBase:
    def __init__(self) -> None:
        self.url = config.api_url
        self.client = httpx.AsyncClient(base_url=self.url)

    async def request(
        self, method: str, path: str, raw: bool, **kwargs
    ) -> Response | dict | Any:
        response = await self.client.request(method, path, **kwargs)
        if raw:
            return response
        return response.json()

    async def post(
        self,
        path: str,
        json: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Response | dict | Any:
        return await self.request("POST", path, json=json, raw=raw, **kwargs)

    async def get(
        self,
        path: str,
        params: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Response | dict | Any:
        return await self.request("GET", path, params=params, raw=raw, **kwargs)

    async def put(
        self,
        path: str,
        data: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Response | dict | Any:
        return await self.request("PUT", path, data=data, raw=raw, **kwargs)

    async def delete(
        self,
        path: str,
        params: Any = None,
        raw: bool = False,
        **kwargs,
    ) -> Response | dict | Any:
        return await self.request("DELETE", path, params=params, raw=raw, **kwargs)
