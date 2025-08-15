from . import NotFoundError


class SystemNotFound(NotFoundError):
    message = "System not found"
