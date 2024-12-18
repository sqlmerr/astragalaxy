use rand::{distributions::Alphanumeric, Rng};
use reqwest::{Client, Method, RequestBuilder};

pub fn request(url: String, method: Method) -> RequestBuilder {
    let client = Client::new();
    client.request(method, url)
}

pub fn snake_to_camel_case(s: &str) -> String {
    let mut result = String::new();
    let mut capitalize_next = false;

    for c in s.chars() {
        match c {
            '_' => capitalize_next = true,
            _ => {
                if capitalize_next {
                    result.push(c.to_uppercase().next().unwrap());
                    capitalize_next = false;
                } else {
                    result.push(c);
                }
            }
        }
    }

    result
}

pub fn snake_case_to_normal(s: &str) -> String {
    s.replace("_", " ")
}

pub fn generate_token(len: usize) -> String {
    rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(len)
        .map(char::from)
        .collect()
}

pub fn parse_token(input: String) -> Option<(i64, String)> {
    let parts = input.split_once(':')?;

    let telegram_id = parts.0.parse::<i64>().ok()?;
    let token = parts.1.to_string();

    Some((telegram_id, token))
}

#[cfg(test)]
mod tests {
    use crate::{generate_token, parse_token};

    #[test]
    fn test_generate_token() {
        let token = generate_token(48);
        println!("{}", token);
        assert_eq!(token.len(), 48);
    }

    #[test]
    fn test_parse_token() {
        let token = generate_token(32);
        let user_id: i64 = 123456789;
        let input_token = format!("{user_id}:{token}");

        let parsed_token = parse_token(input_token);
        assert!(parsed_token.is_some());
        let t = parsed_token.unwrap();
        assert_eq!(t.0, user_id);
        assert_eq!(t.1, token);
    }
}
