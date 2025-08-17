from dataclasses import dataclass

from . import NotFoundError, AppError


@dataclass(eq=False)
class SpaceshipNotFoundError(NotFoundError):
    @property
    def message(self) -> str:
        return "Spaceship not found"


@dataclass(eq=False)
class TooManySpaceshipsError(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Character have too many spaceships"


@dataclass(eq=False)
class InvalidSpaceshipNameError(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Invalid spaceship name. Must be at least 3 characters length and maximum 32 characters."


@dataclass(eq=False)
class CharacterAlreadyInSpaceship(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Character already in spaceship"


@dataclass(eq=False)
class CharacterAlreadyOutOfSpaceship(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Character already out of spaceship"


@dataclass(eq=False)
class CharacterNeedsToBeInSpaceship(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Character needs to be in spaceship"


@dataclass(eq=False)
class InvalidHyperjumpPath(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "Invalid hyperjump path"
