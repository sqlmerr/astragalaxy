import random
import string


def generate_user_token(length: int = 32) -> str:
    out = ""
    for _ in range(length):
        out += random.choice(string.ascii_letters + string.digits + string.punctuation)

    return out


def generate_random_id(length: int = 16) -> str:
    out = ""
    for _ in range(length):
        out += random.choice(string.ascii_letters)

    return out
