from . import NotFoundError, AppError


class SpaceshipNotFoundError(NotFoundError):
    @property
    def message(self) -> str:
        return "Spaceship not found"


class TooManySpaceshipsError(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Character have too many spaceships"


class InvalidSpaceshipNameError(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Invalid spaceship name. Must be at least 3 characters length and maximum 32 characters."
