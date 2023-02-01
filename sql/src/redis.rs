use anyhow::Result;
use redis::Client;

pub async fn open(password: String, url: String) -> Result<Client> {
    let url = format!("redis://:{password}@{url}");
    let client = Client::open(url)?;

    Ok(client)
}
