use anyhow::Result;
use axum::extract::Path;
use axum::http::StatusCode;
use axum::{Extension, Json};
use redis::{Client, Commands};
use serde::{Deserialize, Serialize};
use serde_json::{from_str, to_string};
use sqlx::{query, query_as, FromRow};
use tracing::info;
use uuid::Uuid;

use crate::error::ResponseError;
#[cfg(feature = "mysql")]
use crate::mysql::{bind_key, DatabasePool};
#[cfg(feature = "postgres")]
use crate::postgres::{bind_key, DatabasePool};
use crate::redis::RedisEnabled;

#[derive(Deserialize)]
pub struct CreateMessage {
    id: Option<Uuid>,
    content: String,
    user_id: Uuid,
}

#[derive(Serialize, Deserialize, FromRow)]
pub struct Message {
    id: Uuid,
    content: String,
    user_id: Uuid,
}

pub async fn create_message(
    Extension(pool): Extension<DatabasePool>,
    Json(payload): Json<CreateMessage>,
) -> Result<(StatusCode, Json<Message>), ResponseError> {
    let mut conn = pool.acquire().await?;
    let message = Message {
        id: payload.id.unwrap_or(Uuid::new_v4()),
        content: payload.content,
        user_id: payload.user_id,
    };
    let insert = format!(
        "insert into messages(id, content, user_id) values({}, {}, {})",
        bind_key(1),
        bind_key(2),
        bind_key(3),
    );
    query(&insert)
        .bind(message.id)
        .bind(&message.content)
        .bind(message.user_id)
        .execute(&mut conn)
        .await?;

    info!(
        "successfully created message {} with id {} for user {}",
        message.content, message.id, message.user_id
    );

    Ok((StatusCode::CREATED, Json(message)))
}

pub async fn delete_message(
    Extension(pool): Extension<DatabasePool>,
    Path(id): Path<String>,
) -> Result<StatusCode, ResponseError> {
    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().await?;
    let delete = format!("delete from messages where id = {}", bind_key(1));
    query(&delete).bind(id).execute(&mut conn).await?;

    info!("successfully deleted message with id {}", id);

    Ok(StatusCode::OK)
}

pub async fn get_message(
    Extension(pool): Extension<DatabasePool>,
    Path(id): Path<String>,
) -> Result<(StatusCode, Json<Message>), ResponseError> {
    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().await?;
    let select = format!("select * from messages where id = {}", bind_key(1));
    let message = query_as::<_, Message>(&select)
        .bind(id)
        .fetch_one(&mut conn)
        .await?;

    info!(
        "successfully got message {} with id {} for user {}",
        message.content, message.id, message.user_id
    );

    Ok((StatusCode::OK, Json(message)))
}

pub async fn get_user_messages(
    Extension(pool): Extension<DatabasePool>,
    Extension(client): Extension<Client>,
    Path(id): Path<String>,
    RedisEnabled(redis_enabled): RedisEnabled,
) -> Result<(StatusCode, Json<Vec<Message>>), ResponseError> {
    if redis_enabled {
        let mut conn = client.get_connection()?;
        let messages: Option<String> = conn.get(id.to_string())?;
        match messages {
            None => {}
            Some(messages) => {
                info!(
                    "successfully got cached messages from redis for user {}",
                    id
                );
                let messages = from_str(&messages)?;
                return Ok((StatusCode::OK, Json(messages)));
            }
        }
    }

    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().await?;
    let select = format!("select * from messages where user_id = {}", bind_key(1));
    let messages = query_as::<_, Message>(&select)
        .bind(id)
        .fetch_all(&mut conn)
        .await?;

    info!(
        "successfully got {} messages from database for user {}",
        messages.len(),
        id
    );

    if redis_enabled {
        let mut conn = client.get_connection()?;
        let messages = to_string(&messages)?;
        conn.set(id.to_string(), messages)?;
        conn.expire(id.to_string(), 60)?;
    }

    Ok((StatusCode::OK, Json(messages)))
}
