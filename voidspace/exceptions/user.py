from . import NotFoundError, AppError


class UserNotFound(NotFoundError):
    message = "User not found"


class UsernameAlreadyOccupied(AppError):
    message = "This username already occupied"
    status = 409
