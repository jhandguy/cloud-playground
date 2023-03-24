use std::net::SocketAddr;

use anyhow::Result;
use axum::routing::{delete, get, post};
use axum::{Extension, Router, Server};
use clap::Parser;
use tower_http::validate_request::ValidateRequestHeaderLayer;
use tracing::info;

use sql::message::{create_message, delete_message, get_message, get_user_messages};
use sql::user::{create_user, delete_user, get_user};

use sql::metrics::{serve_metrics, MetricsLayer};
#[cfg(feature = "mysql")]
use sql::mysql::{connect, migrate};
#[cfg(feature = "postgres")]
use sql::postgres::{connect, migrate};
use sql::redis::open;
use sql::tracing::{start_tracing, stop_tracing};

#[derive(Parser)]
pub struct Args {
    /// Service port
    #[clap(long, env)]
    pub sql_http_port: u16,

    /// Metrics port
    #[clap(long, env)]
    pub sql_metrics_port: u16,

    /// Service token
    #[clap(long, env)]
    pub sql_token: String,

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

    /// Redis URL
    #[clap(long, env)]
    pub redis_url: String,

    /// Redis password
    #[clap(long, env)]
    pub redis_password: String,

    /// Tempo URL
    #[clap(long, env)]
    pub tempo_url: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();
    start_tracing(args.tempo_url)?;
    let metrics = serve_metrics("/monitoring/metrics", args.sql_metrics_port).await?;

    info!("connecting to database");
    let pool = connect(
        args.database_user,
        args.database_password,
        args.database_url,
        args.database_name,
    )
    .await?;

    info!("starting data migration");
    migrate(&pool).await?;

    let client = open(args.redis_password, args.redis_url).await?;
    let router = Router::new()
        .route("/message", post(create_message))
        .route("/message/:id", get(get_message))
        .route("/message/:id", delete(delete_message))
        .route("/user", post(create_user))
        .route("/user/:id", get(get_user))
        .route("/user/:id", delete(delete_user))
        .route("/user/:id/messages", get(get_user_messages))
        .route_layer(ValidateRequestHeaderLayer::bearer(&args.sql_token))
        .layer(Extension(pool))
        .layer(Extension(client))
        .layer(MetricsLayer(metrics));

    info!("listening on http port {}", args.sql_http_port);
    let addr = SocketAddr::from(([0, 0, 0, 0], args.sql_http_port));
    Server::bind(&addr)
        .serve(router.into_make_service())
        .await?;

    stop_tracing();

    Ok(())
}
