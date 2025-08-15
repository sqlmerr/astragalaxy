from fastapi import Depends
from fastapi.security import HTTPBearer, APIKeyHeader

from voidspace.api.schemas import Pagination

bearer_security = HTTPBearer()
character_security = APIKeyHeader(
    name="X-Character-ID", description="Active character id"
)


def pagination_params(page: int = 0, per_page: int = 10):
    return Pagination(page=page, per_page=per_page)


JwtSecurity = Depends(bearer_security)
CharacterSecurity = Depends(character_security)
PaginationDepends = Depends(pagination_params)
