pub fn validate_string(s: String) -> bool {
    let len = s.len();
    if len <= 3 || len >= 32 || s.contains(' ') {
        return false;
    }

    let allowed_characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ\
                               абвгдеёжзиклмнопрстуфхцчшщъыьэюя-_@#$%!*0123456789";

    if s.chars().all(|c| c.is_digit(10)) {
        return false;
    }

    s.chars().all(|c| allowed_characters.contains(c))
}
