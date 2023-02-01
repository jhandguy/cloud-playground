use anyhow::Result;
use axum::http::StatusCode;
use axum::{extract, Extension, Json};
use redis::{Client, Commands};
use serde::{Deserialize, Serialize};
use sqlx::{query, query_as, FromRow};
use uuid::Uuid;

use crate::error::ResponseError;
#[cfg(feature = "mysql")]
use crate::mysql::{bind_key, DatabasePool};
#[cfg(feature = "postgres")]
use crate::postgres::{bind_key, DatabasePool};

#[derive(Deserialize)]
pub struct CreateMessage {
    content: String,
}

#[derive(Serialize, FromRow)]
pub struct Message {
    id: Uuid,
    content: String,
}

pub async fn create_message(
    Extension(pool): Extension<DatabasePool>,
    Json(payload): Json<CreateMessage>,
) -> Result<(StatusCode, Json<Message>), ResponseError> {
    let mut conn = pool.acquire().await?;
    let message = Message {
        id: Uuid::new_v4(),
        content: payload.content,
    };
    let insert = format!(
        "insert into message(id, content) values({}, {})",
        bind_key(1),
        bind_key(2)
    );
    query(&insert)
        .bind(message.id)
        .bind(&message.content)
        .execute(&mut conn)
        .await?;

    Ok((StatusCode::CREATED, Json(message)))
}

pub async fn delete_message(
    Extension(pool): Extension<DatabasePool>,
    Extension(client): Extension<Client>,
    extract::Path(id): extract::Path<String>,
) -> Result<StatusCode, ResponseError> {
    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().await?;
    let delete = format!("delete from message where id = {}", bind_key(1));
    query(&delete).bind(id).execute(&mut conn).await?;

    let mut redis = client.get_connection()?;
    redis.del(id.to_string())?;

    Ok(StatusCode::OK)
}

pub async fn get_message(
    Extension(pool): Extension<DatabasePool>,
    Extension(client): Extension<Client>,
    extract::Path(id): extract::Path<String>,
) -> Result<(StatusCode, Json<Message>), ResponseError> {
    let id = Uuid::parse_str(&id)?;
    let mut redis = client.get_connection()?;
    let content: Option<String> = redis.get(id.to_string())?;

    match content {
        None => {}
        Some(content) => {
            return Ok((StatusCode::OK, Json(Message { id, content })));
        }
    }

    let mut conn = pool.acquire().await?;
    let select = format!("select * from message where id = {}", bind_key(1));
    let message = query_as::<_, Message>(&select)
        .bind(id)
        .fetch_one(&mut conn)
        .await?;

    redis.set(id.to_string(), &message.content)?;
    redis.expire(id.to_string(), 60)?;

    Ok((StatusCode::OK, Json(message)))
}
