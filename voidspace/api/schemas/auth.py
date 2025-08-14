from pydantic import BaseModel

from voidspace.dto.user import CreateUserDTO, LoginUserDTO


class AuthRegisterSchema(BaseModel):
    username: str
    password: str

    def into_dto(self) -> CreateUserDTO:
        return CreateUserDTO(username=self.username, password=self.password)


class AuthLoginSchema(BaseModel):
    username: str
    password: str

    def into_dto(self) -> LoginUserDTO:
        return LoginUserDTO(username=self.username, password=self.password)


class TokenSchema(BaseModel):
    access_token: str
    token_type: str = "Bearer"
