use std::str::FromStr;
use std::time::Duration;

use anyhow::Result;
use sqlx::mysql::{MySqlConnectOptions, MySqlPoolOptions};
use sqlx::{migrate, ConnectOptions, MySqlPool};

pub type DatabasePool = MySqlPool;

pub async fn connect(
    user: String,
    password: String,
    url: String,
    database: String,
) -> Result<DatabasePool> {
    let url = format!("mysql://{user}:{password}@{url}/{database}");
    let options = MySqlConnectOptions::from_str(&url)?
        .disable_statement_logging()
        .clone();
    let pool = MySqlPoolOptions::new()
        .acquire_timeout(Duration::from_secs(3))
        .max_connections(5)
        .connect_with(options)
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
