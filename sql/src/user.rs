use anyhow::Result;
use axum::extract::Path;
use axum::http::StatusCode;
use axum::{Extension, Json};
use opentelemetry::global::tracer;
use opentelemetry::trace::{FutureExt, TraceContextExt, Tracer};
use opentelemetry::Context;
use redis::{Client, Commands};
use serde::{Deserialize, Serialize};
use serde_json::to_string;
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
pub struct CreateUser {
    id: Option<Uuid>,
    name: String,
}

#[derive(Serialize, FromRow)]
pub struct User {
    id: Uuid,
    name: String,
}

pub async fn create_user(
    Extension(pool): Extension<DatabasePool>,
    Json(payload): Json<CreateUser>,
) -> Result<(StatusCode, Json<User>), ResponseError> {
    let tracer = tracer("user/create_user");
    let span = tracer.start("user/create_user");
    let ctx = Context::current_with_span(span);

    let mut conn = pool.acquire().with_context(ctx.clone()).await?;
    let user = User {
        id: payload.id.unwrap_or(Uuid::new_v4()),
        name: payload.name,
    };
    let insert = format!(
        "insert into users(id, name) values({}, {})",
        bind_key(1),
        bind_key(2),
    );
    query(&insert)
        .bind(user.id)
        .bind(&user.name)
        .execute(&mut conn)
        .with_context(ctx.clone())
        .await?;

    info!(
        user = to_string(&user)?,
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully created user"
    );

    Ok((StatusCode::CREATED, Json(user)))
}

pub async fn delete_user(
    Extension(pool): Extension<DatabasePool>,
    Extension(client): Extension<Client>,
    Path(id): Path<String>,
    RedisEnabled(redis_enabled): RedisEnabled,
) -> Result<StatusCode, ResponseError> {
    let tracer = tracer("user/delete_user");
    let span = tracer.start("user/delete_user");
    let ctx = Context::current_with_span(span);

    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().with_context(ctx.clone()).await?;
    let delete = format!("delete from messages where user_id = {}", bind_key(1));
    query(&delete)
        .bind(id)
        .execute(&mut conn)
        .with_context(ctx.clone())
        .await?;

    let delete = format!("delete from users where id = {}", bind_key(1));
    query(&delete)
        .bind(id)
        .execute(&mut conn)
        .with_context(ctx.clone())
        .await?;

    if redis_enabled {
        let mut conn = client.get_connection()?;
        conn.del(id.to_string())?;
    }

    info!(
        id = id.to_string(),
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully deleted user"
    );

    Ok(StatusCode::OK)
}

pub async fn get_user(
    Extension(pool): Extension<DatabasePool>,
    Path(id): Path<String>,
) -> Result<(StatusCode, Json<User>), ResponseError> {
    let tracer = tracer("user/get_user");
    let span = tracer.start("user/get_user");
    let ctx = Context::current_with_span(span);

    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().with_context(ctx.clone()).await?;
    let select = format!("select * from users where id = {}", bind_key(1));
    let user = query_as::<_, User>(&select)
        .bind(id)
        .fetch_one(&mut conn)
        .with_context(ctx.clone())
        .await?;

    info!(
        user = to_string(&user)?,
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully got user"
    );

    Ok((StatusCode::OK, Json(user)))
}
