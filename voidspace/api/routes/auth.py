from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from voidspace.api.dependencies import JwtSecurity
from voidspace.api.schemas.auth import AuthRegisterSchema, AuthLoginSchema, TokenSchema
from voidspace.api.schemas.user import UserSchema
from voidspace.identity_provider import IdentityProvider
from voidspace.interfaces.user import UserWriter, UserReader

router = APIRouter(prefix="/auth", route_class=DishkaRoute, tags=["Auth"])


@router.post("/register")
async def register_user(
    data: AuthRegisterSchema, user_writer: FromDishka[UserWriter]
) -> UserSchema:
    user = await user_writer.create_user(data.into_dto())

    return UserSchema.from_dto(user)


@router.post("/login")
async def login(
    data: AuthLoginSchema, user_reader: FromDishka[UserReader]
) -> TokenSchema:
    token = await user_reader.login(data.into_dto())
    return TokenSchema(access_token=token)


@router.get("/me", dependencies=[JwtSecurity])
async def get_me(identity_provider: FromDishka[IdentityProvider]) -> UserSchema:
    user = await identity_provider.get_current_user()
    return UserSchema.from_dto(user)
