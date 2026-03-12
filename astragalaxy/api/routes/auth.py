from dishka import FromDishka
from dishka.integrations.fastapi import DishkaRoute
from fastapi import APIRouter

from astragalaxy.api.dependencies import JwtSecurity
from astragalaxy.api.schemas.auth import (
    AuthRegisterSchema,
    AuthLoginSchema,
    TokenSchema,
)
from astragalaxy.api.schemas.user import UserSchema
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.use_cases.login import Login
from astragalaxy.use_cases.register import Register

router = APIRouter(prefix="/auth", route_class=DishkaRoute, tags=["Auth"])


@router.post("/register")
async def register_user(
    data: AuthRegisterSchema, use_case: FromDishka[Register]
) -> UserSchema:
    user = await use_case.execute(data.into_dto())
    return UserSchema.from_dto(user)


@router.post("/login")
async def login(data: AuthLoginSchema, use_case: FromDishka[Login]) -> TokenSchema:
    token = await use_case.execute(data.into_dto())
    return TokenSchema(access_token=token.access_token, token_type=token.token_type)


@router.get("/me", dependencies=[JwtSecurity])
async def get_me(identity_provider: FromDishka[IdentityProvider]) -> UserSchema:
    user = await identity_provider.get_current_user()
    return UserSchema.from_dto(user)
