use anyhow::Result;
use axum::extract::Path;
use axum::http::StatusCode;
use axum::{Extension, Json};
use redis::{Client, Commands};
use serde::{Deserialize, Serialize};
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
    let mut conn = pool.acquire().await?;
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
        .await?;

    info!(
        "successfully created user {} with id {}",
        user.name, user.id
    );

    Ok((StatusCode::CREATED, Json(user)))
}

pub async fn delete_user(
    Extension(pool): Extension<DatabasePool>,
    Extension(client): Extension<Client>,
    Path(id): Path<String>,
    RedisEnabled(redis_enabled): RedisEnabled,
) -> Result<StatusCode, ResponseError> {
    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().await?;
    let delete = format!("delete from messages where user_id = {}", bind_key(1));
    query(&delete).bind(id).execute(&mut conn).await?;

    let delete = format!("delete from users where id = {}", bind_key(1));
    query(&delete).bind(id).execute(&mut conn).await?;

    if redis_enabled {
        let mut conn = client.get_connection()?;
        conn.del(id.to_string())?;
    }

    info!("successfully deleted user with id {}", id);

    Ok(StatusCode::OK)
}

pub async fn get_user(
    Extension(pool): Extension<DatabasePool>,
    Path(id): Path<String>,
) -> Result<(StatusCode, Json<User>), ResponseError> {
    let id = Uuid::parse_str(&id)?;
    let mut conn = pool.acquire().await?;
    let select = format!("select * from users where id = {}", bind_key(1));
    let user = query_as::<_, User>(&select)
        .bind(id)
        .fetch_one(&mut conn)
        .await?;

    info!("successfully got user {} with id {}", user.name, user.id);

    Ok((StatusCode::OK, Json(user)))
}
