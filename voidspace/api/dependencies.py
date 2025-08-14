from fastapi import Depends
from fastapi.security import HTTPBearer, APIKeyHeader


bearer_security = HTTPBearer()
character_security = APIKeyHeader(
    name="X-Character-ID", description="Active character id"
)

JwtSecurity = Depends(bearer_security)
CharacterSecurity = Depends(character_security)
