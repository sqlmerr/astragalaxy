def validate_string(s: str) -> bool:
    if 3 < len(s) < 32 and " " not in s:
        allowed_characters = set(
            "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
            "абвгдеёжзиклмнопрстуфхцчшщъыьэюя"
            "-_@#$%!*0123456789"
        )

        if s.isdigit():
            return False

        return all(char in allowed_characters for char in s)

    return False
