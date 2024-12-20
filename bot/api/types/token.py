from pydantic import BaseModel


class TokenPair(BaseModel):
    user_token: str
    jwt_token: str
