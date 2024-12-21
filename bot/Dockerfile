FROM python:3.12-slim

ENV PYTHONUNBUFFERED=1 \
    # prevents python creating .pyc as files
    PYTHONDONTWRITEBYTECODE=1 \
    \
    # pip
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on \
    PIP_DEFAULT_TIMEOUT=100 \
    \
    # poetry
    # https://python-poetry.org/docs/configuration/#using-environment-variables
    POETRY_VERSION=1.6.1 \
    # make poetry install to this location
    POETRY_HOME="/opt/poetry" \
    # make poetry create the virtual environment in the project's root
    # it gets named `.venv`
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    # do not ask any interactive question
    POETRY_NO_INTERACTION=1 \
    \
    # paths
    # this is where our requirements + virtual environment will live
    PYSETUP_PATH="/opt/pysetup" \
    VENV_PATH="/opt/pysetup/.venv"

WORKDIR /app

COPY poetry.lock pyproject.toml ./

RUN python -m pip install --no-cache-dir poetry==1.6.1 \
    && poetry config virtualenvs.create false \
    && poetry install --without dev --no-interaction --no-ansi \
    && rm -rf $(poetry config cache-dir)/{cache,artifacts}

RUN poetry install --no-dev

COPY . .

CMD ["python", "main.py"]