use crate::error::ResponseError;
use anyhow::Result;
use axum::extract::FromRequestParts;
use axum::http::request::Parts;
use axum::http::HeaderMap;
use axum::{async_trait, RequestPartsExt};
use redis::Client;

pub async fn open(password: String, url: String) -> Result<Client> {
    let url = format!("redis://:{password}@{url}");
    let client = Client::open(url)?;

    Ok(client)
}

pub struct RedisEnabled(pub bool);

#[async_trait]
impl<S> FromRequestParts<S> for RedisEnabled
where
    S: Send + Sync,
{
    type Rejection = ResponseError;

    async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
        let headers = parts.extract::<HeaderMap>().await?;
        let header = headers.get("X-Redis-Enabled");

        return match header {
            None => Ok(RedisEnabled(false)),
            Some(value) => Ok(RedisEnabled(value.eq("true"))),
        };
    }
}
