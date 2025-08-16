from dataclasses import dataclass

from . import NotFoundError, AppError


@dataclass(eq=False)
class UserNotFound(NotFoundError):
    @property
    def message(self) -> str:
        return "User not found"


@dataclass(eq=False)
class UsernameAlreadyOccupied(AppError):
    @property
    def status(self) -> int:
        return 409

    @property
    def message(self) -> str:
        return "This username already occupied"


@dataclass(eq=False)
class InvalidCredentials(AppError):
    @property
    def status(self) -> int:
        return 401

    @property
    def message(self) -> str:
        return "Invalid credentials"


class InvalidToken(AppError):
    @property
    def status(self) -> int:
        return 401

    @property
    def message(self) -> str:
        return "Invalid token"
