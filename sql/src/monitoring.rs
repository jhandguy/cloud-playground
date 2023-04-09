use anyhow::Result;
use axum::http::StatusCode;
use axum::Extension;
use opentelemetry::global::tracer;
use opentelemetry::trace::{FutureExt, TraceContextExt, Tracer};
use opentelemetry::Context;
use sqlx::Connection;
use tracing::debug;

use crate::error::ResponseError;
#[cfg(feature = "mysql")]
use crate::mysql::DatabasePool;
#[cfg(feature = "postgres")]
use crate::postgres::DatabasePool;

pub async fn check_readiness(
    Extension(pool): Extension<DatabasePool>,
) -> Result<StatusCode, ResponseError> {
    let tracer = tracer("monitoring/check_readiness");
    let span = tracer.start("monitoring/check_readiness");
    let ctx = Context::current_with_span(span);

    let mut conn = pool.acquire().with_context(ctx.clone()).await?;
    conn.ping().await?;

    debug!(
        trace_id = ctx.span().span_context().trace_id().to_string(),
        "successfully checked readiness"
    );

    Ok(StatusCode::OK)
}

pub async fn check_liveness() -> Result<StatusCode, ResponseError> {
    Ok(StatusCode::OK)
}
