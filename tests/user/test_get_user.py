from unittest.mock import Mock, AsyncMock
from uuid import uuid4

import pytest

from astragalaxy.database.models import User
from astragalaxy.dto.user import UserDTO
from astragalaxy.exceptions import AccessDeniedError
from astragalaxy.identity_provider import IdentityProvider
from astragalaxy.interfaces.user import UserRepo
from astragalaxy.use_cases.get_user import GetUserById

pytestmark = pytest.mark.asyncio

ADMIN_USER_ID = uuid4()
SECOND_USER_ID = uuid4()


@pytest.fixture
def idp() -> IdentityProvider:
    identity_provider = Mock()
    identity_provider.get_current_user_id = Mock(return_value=ADMIN_USER_ID)
    identity_provider.get_current_user = AsyncMock(
        return_value=UserDTO(
            id=ADMIN_USER_ID, username="admin", password="password", token="token"
        )
    )
    identity_provider.get_current_character_id = Mock()
    identity_provider.get_current_character = AsyncMock()

    return identity_provider


@pytest.fixture
def idp_2() -> IdentityProvider:
    identity_provider = Mock()
    identity_provider.get_current_user_id = Mock(return_value=SECOND_USER_ID)
    identity_provider.get_current_user = AsyncMock(
        return_value=UserDTO(
            id=SECOND_USER_ID,
            username="second_user",
            password="password",
            token="token",
        )
    )
    identity_provider.get_current_character_id = Mock()
    identity_provider.get_current_character = AsyncMock()

    return identity_provider


@pytest.fixture
def user_repo() -> UserRepo:
    repo = Mock()
    repo.find_one_user = AsyncMock(
        return_value=User(
            id=SECOND_USER_ID,
            username="second_user",
            password="password",
            token="token",
        )
    )
    repo.find_one_user_by_username = AsyncMock(
        return_value=User(
            id=SECOND_USER_ID,
            username="second_user",
            password="password",
            token="token",
        )
    )

    return repo


@pytest.fixture
def get_user_by_id(idp: IdentityProvider, user_repo: UserRepo) -> GetUserById:
    return GetUserById(user_repo, idp)


@pytest.fixture
def get_user_by_id_2(idp_2: IdentityProvider, user_repo: UserRepo) -> GetUserById:
    return GetUserById(user_repo, idp_2)


async def test_get_user_by_id(
    get_user_by_id: GetUserById, get_user_by_id_2: GetUserById
):
    with pytest.raises(AccessDeniedError):
        await get_user_by_id.execute(SECOND_USER_ID)

    result = await get_user_by_id_2.execute(SECOND_USER_ID)
    assert result.id == SECOND_USER_ID
