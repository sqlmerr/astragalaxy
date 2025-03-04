from . import planet, hyperjump, info

class Flights(planet.PlanetFlight, hyperjump.HyperJump, info.GetFlightInfo):
    pass