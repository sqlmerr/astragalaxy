from . import register, login, get_token, get_me


class Auth(
    register.Register,
    login.Login,
    get_token.GetToken,
    get_me.GetMe,
):
    pass
