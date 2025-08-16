from dataclasses import dataclass


@dataclass(eq=False)
class AppError(Exception):
    @property
    def status(self) -> int:
        return 500

    @property
    def message(self) -> str:
        return "Server Error"


@dataclass(eq=False)
class NotFoundError(AppError):
    @property
    def status(self) -> int:
        return 404

    @property
    def message(self) -> str:
        return "Not found"


@dataclass(eq=False)
class AccessDeniedError(AppError):
    @property
    def status(self) -> int:
        return 403

    @property
    def message(self) -> str:
        return "Access denied"
