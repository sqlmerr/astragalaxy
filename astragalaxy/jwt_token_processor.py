from typing import Any

import jwt

from astragalaxy.config import Settings


class JwtTokenProcessor:
    def __init__(self, config: Settings):
        self.config = config

    def encode(self, payload: dict[str, Any]):
        return jwt.encode(payload, self.config.jwt_secret, algorithm="HS256")

    def decode(self, token: str | bytes):
        return jwt.decode(token, self.config.jwt_secret, algorithms=["HS256"])
