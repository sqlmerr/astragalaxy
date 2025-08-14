from . import NotFoundError, AppError


class UserNotFound(NotFoundError):
    message = "User not found"


class UsernameAlreadyOccupied(AppError):
    message = "This username already occupied"
    status = 409


class InvalidCredentials(AppError):
    status = 401
    message = "Invalid credentials"


class InvalidToken(AppError):
    status = 401
    message = "Invalid token"
