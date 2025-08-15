from pydantic import BaseModel


class DataSchema[T: BaseModel](BaseModel):
    data: list[T]


class Pagination(BaseModel):
    per_page: int = 10
    page: int = 0
