use anyhow::Result;
use axum::extract::Path;
use axum::http::StatusCode;
use axum::{Extension, Json};
use opentelemetry::global::tracer;
use opentelemetry::trace::{FutureExt, TraceContextExt, Tracer};
use opentelemetry::Context;
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
    let tracer = tracer("message/create_message");
    let span = tracer.start("message/create_message");
    let ctx = Context::current_with_span(span);
    let mut tx = pool.begin().with_context(ctx.clone()).await?;

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
        .execute(tx.as_mut())
        .with_context(ctx.clone())
        .await?;

    tx.commit().await?;

    info!(
        msg = to_string(&message)?,
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully created message"
    );

    Ok((StatusCode::CREATED, Json(message)))
}

pub async fn delete_message(
    Extension(pool): Extension<DatabasePool>,
    Path(id): Path<String>,
) -> Result<StatusCode, ResponseError> {
    let tracer = tracer("message/delete_message");
    let span = tracer.start("message/delete_message");
    let ctx = Context::current_with_span(span);
    let mut tx = pool.begin().with_context(ctx.clone()).await?;

    let id = Uuid::parse_str(&id)?;
    let delete = format!("delete from messages where id = {}", bind_key(1));
    query(&delete)
        .bind(id)
        .execute(tx.as_mut())
        .with_context(ctx.clone())
        .await?;

    tx.commit().await?;

    info!(
        id = id.to_string(),
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully deleted message"
    );

    Ok(StatusCode::OK)
}

pub async fn get_message(
    Extension(pool): Extension<DatabasePool>,
    Path(id): Path<String>,
) -> Result<(StatusCode, Json<Message>), ResponseError> {
    let tracer = tracer("message/get_message");
    let span = tracer.start("message/get_message");
    let ctx = Context::current_with_span(span);
    let mut conn = pool.acquire().with_context(ctx.clone()).await?;

    let id = Uuid::parse_str(&id)?;
    let select = format!("select * from messages where id = {}", bind_key(1));
    let message = query_as::<_, Message>(&select)
        .bind(id)
        .fetch_one(conn.as_mut())
        .with_context(ctx.clone())
        .await?;

    info!(
        msg = to_string(&message)?,
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully got message"
    );

    Ok((StatusCode::OK, Json(message)))
}

pub async fn get_user_messages(
    Extension(pool): Extension<DatabasePool>,
    Extension(client): Extension<Client>,
    Path(id): Path<String>,
    RedisEnabled(redis_enabled): RedisEnabled,
) -> Result<(StatusCode, Json<Vec<Message>>), ResponseError> {
    let tracer = tracer("message/get_user_messages");
    let span = tracer.start("message/get_user_messages");
    let ctx = Context::current_with_span(span);

    if redis_enabled {
        let mut conn = client.get_connection()?;
        let messages: Option<String> = conn.get(id.to_string())?;
        match messages {
            None => {}
            Some(messages) => {
                info!(
                    user_id = id,
                    trace_id = ctx.span().span_context().trace_id().to_string(),
                    "successfully got messages from redis"
                );
                let messages = from_str(&messages)?;
                return Ok((StatusCode::OK, Json(messages)));
            }
        }
    }

    let mut conn = pool.acquire().with_context(ctx.clone()).await?;
    let id = Uuid::parse_str(&id)?;
    let select = format!("select * from messages where user_id = {}", bind_key(1));
    let messages = query_as::<_, Message>(&select)
        .bind(id)
        .fetch_all(conn.as_mut())
        .with_context(ctx.clone())
        .await?;

    info!(
        user_id = id.to_string(),
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully got messages from database"
    );

    if redis_enabled {
        let mut conn = client.get_connection()?;
        let messages = to_string(&messages)?;
        conn.set(id.to_string(), messages)?;
        conn.expire(id.to_string(), 60)?;
    }

    Ok((StatusCode::OK, Json(messages)))
}
