from pydantic import BaseModel


class DataSchema[T: BaseModel](BaseModel):
    data: list[T]
