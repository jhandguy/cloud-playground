use std::net::SocketAddr;

use anyhow::Result;
use axum::routing::{delete, get, post};
use axum::{Extension, Router};
use clap::Parser;

use sql::message::{create_message, delete_message, get_message};

#[cfg(feature = "mysql")]
use sql::mysql::{connect, migrate};
#[cfg(feature = "postgres")]
use sql::postgres::{connect, migrate};
use sql::redis::open;

#[derive(Parser)]
pub struct Args {
    /// Service port
    #[clap(long, env)]
    pub sql_http_port: u16,

    /// Database URL
    #[clap(long, env)]
    pub database_url: String,

    /// Database user
    #[clap(long, env)]
    pub database_user: String,

    /// Database password
    #[clap(long, env)]
    pub database_password: String,

    /// Database name
    #[clap(long, env)]
    pub database_name: String,

    /// Redis url
    #[clap(long, env)]
    pub redis_url: String,

    /// Redis password
    #[clap(long, env)]
    pub redis_password: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();
    let pool = connect(
        args.database_user,
        args.database_password,
        args.database_url,
        args.database_name,
    )
    .await?;
    migrate(&pool).await?;

    let client = open(args.redis_password, args.redis_url).await?;

    let router = Router::new()
        .route("/message", post(create_message))
        .route("/message/:id", get(get_message))
        .route("/message/:id", delete(delete_message))
        .layer(Extension(pool))
        .layer(Extension(client));

    let addr = SocketAddr::from(([0, 0, 0, 0], args.sql_http_port));
    axum::Server::bind(&addr)
        .serve(router.into_make_service())
        .await?;

    Ok(())
}
