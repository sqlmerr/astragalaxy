import abc


class BaseUseCase[I, O](abc.ABC):
    @abc.abstractmethod
    async def execute(self, data: I) -> O:
        raise NotImplementedError
