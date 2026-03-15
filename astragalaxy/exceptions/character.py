from dataclasses import dataclass

from . import NotFoundError, AppError


@dataclass(eq=False)
class CharacterNotFound(NotFoundError):
    @property
    def message(self) -> str:
        return "Character not found"


@dataclass(eq=False)
class InvalidCharacterCode(AppError):
    @property
    def status(self) -> int:
        return 400
    
    @property
    def message(self) -> str:
        return "Invalid character code. Allowed characters: 0-9, a-z, A-Z, _"

@dataclass(eq=False)
class CharacterCodeAlreadyOccupied(AppError):
    @property
    def status(self) -> int:
        return 409

    @property
    def message(self) -> str:
        return "Character code already occupied"


@dataclass(eq=False)
class TooManyCharacters(AppError):
    @property
    def status(self) -> int:
        return 400

    @property
    def message(self) -> str:
        return "You have too many characters"
