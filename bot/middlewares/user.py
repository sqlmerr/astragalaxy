from typing import Callable, Awaitable, Any

from aiogram import BaseMiddleware
from aiogram.types import Message, CallbackQuery
from aiogram_i18n import I18nContext

from api import Api
from api.exceptions import AuthError
from api.types.token import TokenPair
from utils.token_manager import TokenManager


class UserMiddleware(BaseMiddleware):
    async def __call__(
        self,
        handler: Callable[[Message, dict[str, Any]], Awaitable[Any]],
        event: Message | CallbackQuery,
        data: dict[str, Any],
    ) -> Any:
        token_manager: TokenManager = data["token_manager"]
        api: Api = data["api"]
        i18n: I18nContext = data["i18n"]
        user_id = event.from_user.id

        jwt_token = await token_manager.get_jwt_token(user_id)
        if jwt_token is None:
            user_token = await token_manager.get_user_token(user_id)
            if user_token is not None:
                token_pair = await api.login_user(user_id, token=user_token)
                await token_manager.set_jwt_token(user_id, token_pair.jwt_token)
            else:
                token_pair = await api.login_user(user_id)
                await token_manager.set_jwt_token(user_id, token_pair.jwt_token)
                await token_manager.set_user_token(user_id, token_pair.user_token)
        else:
            user_token = await token_manager.get_user_token(user_id)
            if user_token is not None:
                token_pair = TokenPair(user_token=user_token, jwt_token=jwt_token)
            else:
                user_token = await api.get_user_token(user_id)
                if user_token is None:
                    if isinstance(event, CallbackQuery):
                        await event.answer(i18n.unexpected_error_callback())
                    elif isinstance(event, Message):
                        await event.reply(i18n.unexpected_error())
                    # TODO: notify admin function, use it here to notify admin about error
                    return
                token_pair = TokenPair(user_token=user_token, jwt_token=jwt_token)
                await token_manager.set_user_token(user_id, token_pair.user_token)

        try:
            user = await api.get_me(token_pair.jwt_token)
        except AuthError:
            token_pair = await api.login_user(user_id, token=token_pair.user_token)
            await token_manager.set_jwt_token(user_id, token_pair.jwt_token)
            user = await api.get_me(token_pair.jwt_token)

        data["token_pair"] = token_pair
        data["user"] = user

        return await handler(event, data)
