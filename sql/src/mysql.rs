use std::time::Duration;

use anyhow::Result;
use sqlx::mysql::MySqlPoolOptions;
use sqlx::{migrate, MySqlPool};

pub type DatabasePool = MySqlPool;

pub async fn connect(
    user: String,
    password: String,
    url: String,
    database: String,
) -> Result<DatabasePool> {
    let url = format!("mysql://{user}:{password}@{url}/{database}");
    let pool = MySqlPoolOptions::new()
        .acquire_timeout(Duration::from_secs(3))
        .max_connections(5)
        .connect(&url)
        .await?;

    Ok(pool)
}

pub async fn migrate(pool: &DatabasePool) -> Result<()> {
    migrate!("./migrations/mysql").run(pool).await?;

    Ok(())
}

pub fn bind_key(_: u8) -> &'static str {
    "?"
}
