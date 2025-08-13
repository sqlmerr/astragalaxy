from typing import ClassVar


class AppError(Exception):
    status: ClassVar[int]
    message: ClassVar[str]


class NotFoundError(AppError):
    status = 404
    message = "Not found"
