use anyhow::Error;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use tracing::error;

pub struct ResponseError(Error);

impl IntoResponse for ResponseError {
    fn into_response(self) -> Response {
        let err = self.0.to_string();
        error!("internal server error: {}", err);
        (StatusCode::INTERNAL_SERVER_ERROR, err).into_response()
    }
}

impl<E> From<E> for ResponseError
where
    E: Into<Error>,
{
    fn from(err: E) -> Self {
        Self(err.into())
    }
}
