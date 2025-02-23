from .auth import Auth
from .spaceships import Spaceships
from .systems import Systems
from .flights import Flights


class Methods(Auth, Spaceships, Systems, Flights):
    pass
