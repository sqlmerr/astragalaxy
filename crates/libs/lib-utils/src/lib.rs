use reqwest::{Client, Method, RequestBuilder};

pub fn request(url: String, method: Method) -> RequestBuilder {
    let client = Client::new();
    client.request(method, url)
}
