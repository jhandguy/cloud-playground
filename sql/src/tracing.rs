use anyhow::Result;
use opentelemetry::global::shutdown_tracer_provider;
use opentelemetry::runtime::Tokio;
use opentelemetry::sdk::trace::config;
use opentelemetry::sdk::Resource;
use opentelemetry::KeyValue;
use opentelemetry_otlp::{new_exporter, new_pipeline, WithExportConfig};
use opentelemetry_semantic_conventions::resource::SERVICE_NAME as SERVICE_NAME_KEY;
use tracing::subscriber::set_global_default;
use tracing::Level;
use tracing_opentelemetry::layer as otl_layer;
use tracing_subscriber::fmt::layer as log_layer;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::{EnvFilter, Registry};

#[cfg(feature = "mysql")]
const SERVICE_NAME: &str = "sql-mysql";
#[cfg(feature = "postgres")]
const SERVICE_NAME: &str = "sql-postgres";

pub fn start_tracing(endpoint: String) -> Result<()> {
    let env_filter = EnvFilter::builder()
        .with_default_directive(Level::INFO.into())
        .from_env_lossy();
    let log_layer = log_layer().json().flatten_event(true);
    let tracer = new_pipeline()
        .tracing()
        .with_exporter(
            new_exporter()
                .tonic()
                .with_endpoint(format!("http://{}", endpoint)),
        )
        .with_trace_config(config().with_resource(Resource::new(vec![KeyValue::new(
            SERVICE_NAME_KEY,
            SERVICE_NAME,
        )])))
        .install_batch(Tokio)?;
    let otl_layer = otl_layer().with_tracer(tracer);
    let subscriber = Registry::default()
        .with(env_filter)
        .with(log_layer)
        .with(otl_layer);

    set_global_default(subscriber)?;

    Ok(())
}

pub fn stop_tracing() {
    shutdown_tracer_provider()
}
