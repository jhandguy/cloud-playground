use crate::error::ResponseError;
use anyhow::Result;
use axum::http::header::CONTENT_TYPE;
use axum::http::{Method, Request, StatusCode};
use axum::response::Response;
use axum::routing::get;
use axum::{Extension, Router, Server};
use pin_project::pin_project;
use prometheus_client::encoding::text::encode;
use prometheus_client::encoding::EncodeLabelSet;
use prometheus_client::metrics::counter::Counter;
use prometheus_client::metrics::family::Family;
use prometheus_client::metrics::histogram::Histogram;
use prometheus_client::registry::Registry;
use std::future::Future;
use std::net::SocketAddr;
use std::pin::Pin;
use std::sync::Arc;
use std::task::Poll::Ready;
use std::task::{ready, Context, Poll};
use tokio::spawn;
use tokio::time::Instant;
use tower::{Layer, Service};
use tracing::info;

#[derive(Clone, Debug, Hash, PartialEq, Eq, EncodeLabelSet)]
pub struct Labels {
    pub method: String,
    pub success: String,
}

#[derive(Clone)]
pub struct Metrics {
    pub request_counter: Family<Labels, Counter>,
    pub latency_histogram: Family<Labels, Histogram>,
}

impl Metrics {
    pub fn inc_request_counter(&self, method: Method, success: bool) {
        self.request_counter
            .get_or_create(&Labels {
                method: method.to_string(),
                success: success.to_string(),
            })
            .inc();
    }

    pub fn observe_latency_histogram(&self, value: f64, method: Method, success: bool) {
        self.latency_histogram
            .get_or_create(&Labels {
                method: method.to_string(),
                success: success.to_string(),
            })
            .observe(value);
    }
}

async fn metrics_handler(
    Extension(registry): Extension<Arc<Registry>>,
) -> Result<(StatusCode, Response<String>), ResponseError> {
    let mut body = String::new();
    encode(&mut body, &registry)?;

    let response = Response::builder()
        .header(
            CONTENT_TYPE,
            "application/openmetrics-text; version=1.0.0; charset=utf-8",
        )
        .body(body)?;

    Ok((StatusCode::OK, response))
}

pub async fn serve_metrics(path: &str, port: u16) -> Result<Metrics> {
    let mut registry = <Registry>::with_prefix("cloud_playground_sql");
    let metrics = Metrics {
        request_counter: Family::<Labels, Counter>::default(),
        latency_histogram: Family::<Labels, Histogram>::new_with_constructor(|| {
            Histogram::new(
                [
                    0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0,
                ]
                .into_iter(),
            )
        }),
    };

    registry.register(
        "requests_count",
        "Request counter per method",
        metrics.request_counter.clone(),
    );

    registry.register(
        "requests_latency",
        "Request latency histogram per method",
        metrics.latency_histogram.clone(),
    );

    let router = Router::new()
        .route(path, get(metrics_handler))
        .layer(Extension(Arc::new(registry)));

    info!("listening on metrics port {}", port);
    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    spawn(async move {
        Server::bind(&addr)
            .serve(router.into_make_service())
            .await
            .unwrap()
    });

    Ok(metrics)
}

#[derive(Clone)]
pub struct MetricsLayer(pub Metrics);

#[derive(Clone)]
pub struct MetricsService<S> {
    metrics: Metrics,
    service: S,
}

#[pin_project]
pub struct MetricsFuture<F> {
    #[pin]
    inner: F,
    metrics: Metrics,
    start: Instant,
    method: Method,
}

impl<S> Layer<S> for MetricsLayer {
    type Service = MetricsService<S>;

    fn layer(&self, service: S) -> Self::Service {
        MetricsService {
            metrics: self.0.clone(),
            service,
        }
    }
}

impl<S, R, B> Service<Request<R>> for MetricsService<S>
where
    S: Service<Request<R>, Response = Response<B>>,
{
    type Response = Response<B>;
    type Error = S::Error;
    type Future = MetricsFuture<S::Future>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.service.poll_ready(cx)
    }

    fn call(&mut self, request: Request<R>) -> Self::Future {
        let start = Instant::now();
        let method = request.method().clone();
        MetricsFuture {
            inner: self.service.call(request),
            metrics: self.metrics.clone(),
            start,
            method,
        }
    }
}

impl<F, B, E> Future for MetricsFuture<F>
where
    F: Future<Output = Result<Response<B>, E>>,
{
    type Output = Result<Response<B>, E>;

    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output> {
        let this = self.project();
        let response = ready!(this.inner.poll(cx))?;
        let success = response.status().is_success();
        let latency = this.start.elapsed().as_secs_f64();

        this.metrics
            .inc_request_counter(this.method.clone(), success);
        this.metrics
            .observe_latency_histogram(latency, this.method.clone(), success);

        Ready(Ok(response))
    }
}
