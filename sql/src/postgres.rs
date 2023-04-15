use std::str::FromStr;

use anyhow::Result;
use sqlx::postgres::{PgConnectOptions, PgPoolOptions};
use sqlx::{migrate, ConnectOptions, PgPool};

pub type DatabasePool = PgPool;

pub async fn connect(
    user: String,
    password: String,
    url: String,
    database: String,
) -> Result<DatabasePool> {
    let url = format!("postgres://{user}:{password}@{url}/{database}");
    let options = PgConnectOptions::from_str(&url)?
        .disable_statement_logging()
        .clone();
    let pool = PgPoolOptions::new().connect_with(options).await?;

    Ok(pool)
}

pub async fn migrate(pool: &DatabasePool) -> Result<()> {
    migrate!("./migrations/postgres").run(pool).await?;

    Ok(())
}

pub fn bind_key(number: u8) -> String {
    format!("${number}")
}
