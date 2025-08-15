from . import NotFoundError, AppError


class CharacterNotFound(NotFoundError):
    message = "Character not found"


class CharacterCodeAlreadyOccupied(AppError):
    status = 409
    message = "Character code already occupied"


class TooManyCharacters(AppError):
    status = 400
    message = "You have too many characters"
