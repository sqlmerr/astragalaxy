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
