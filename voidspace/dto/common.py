from dataclasses import dataclass


@dataclass(frozen=True)
class PaginationDTO:
    per_page: int
    page: int
