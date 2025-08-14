from . import NotFoundError, AppError


class CharacterNotFound(NotFoundError):
    message = "Character not found"


class CharacterCodeAlreadyOccupied(AppError):
    status = 409
    message = "Character code already occupied"
