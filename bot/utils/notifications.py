import traceback

from aiogram import Bot
from aiogram.types import LinkPreviewOptions
from aiogram_i18n import I18nContext

from aiogram.types import User
from config_reader import config


async def notify_admins_error(
    bot: Bot, exc: Exception, i18n: I18nContext, user: User
) -> None:
    print(traceback.format_exc())
    await notify_admins(
        bot,
        i18n.error_admin_notification(
            user=f"<a href='tg://user?id={user.id}'>{user.username}</a>",
            error=f"{repr(exc)} {str(exc)}",
        ),
        link_preview_options=LinkPreviewOptions(is_disabled=True),
    )


async def notify_admins(bot: Bot, text: str, **kwargs) -> None:
    for admin in config.admins:
        await bot.send_message(admin, text, **kwargs)
