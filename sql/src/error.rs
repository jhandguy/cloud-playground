use anyhow::Error;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};

pub struct ResponseError(Error);

impl IntoResponse for ResponseError {
    fn into_response(self) -> Response {
        (StatusCode::INTERNAL_SERVER_ERROR, self.0.to_string()).into_response()
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
