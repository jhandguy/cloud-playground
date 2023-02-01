use std::time::Duration;

use anyhow::Result;
use sqlx::postgres::PgPoolOptions;
use sqlx::{migrate, PgPool};

pub type DatabasePool = PgPool;

pub async fn connect(
    user: String,
    password: String,
    url: String,
    database: String,
) -> Result<DatabasePool> {
    let url = format!("postgres://{user}:{password}@{url}/{database}");
    let pool = PgPoolOptions::new()
        .acquire_timeout(Duration::from_secs(3))
        .max_connections(5)
        .connect(&url)
        .await?;

    Ok(pool)
}

pub async fn migrate(pool: &DatabasePool) -> Result<()> {
    migrate!("./migrations/postgres").run(pool).await?;

    Ok(())
}

pub fn bind_key(number: u8) -> String {
    format!("${number}")
}
