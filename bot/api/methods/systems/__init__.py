from . import get_all, get_by_id, get_planets


class Systems(get_all.GetAll, get_by_id.GetById, get_planets.GetPlanets):
    pass
