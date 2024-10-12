use axum::{
    body::Bytes,
    extract::{Multipart, State},
    response::IntoResponse,
    Json,
};
use serde::Serialize;
use uuid::Uuid;

use crate::{jwt::jwt::Claims, rabbitmq_client::models::MediaUploadedMessage, AppState};
#[derive(Serialize)]
pub struct Message {
    pub message: String,
}

pub async fn check_health() -> impl IntoResponse {
    Json(Message {
        message: "Hello, Axum!".to_string(),
    })
}

#[derive(Debug)]
pub struct UploadError(String);

impl IntoResponse for UploadError {
    fn into_response(self) -> axum::response::Response {
        (axum::http::StatusCode::BAD_REQUEST, self.0).into_response()
    }
}

pub async fn upload_media(
    State(state): State<AppState>,
    claims: Claims,
    mut multipart: Multipart,
) -> Result<impl IntoResponse, UploadError> {
    println!("running handler, claims: {}", claims);
    let mut description = String::new();
    let mut image_file: Option<Bytes> = None;
    while let Some(field) = multipart
        .next_field()
        .await
        .map_err(|e| UploadError(e.to_string()))?
    {
        let name = field
            .name()
            .ok_or_else(|| UploadError("name not found".to_string()));
        let content_type = field
            .content_type()
            .ok_or_else(|| UploadError("content type not found".to_string()));

        match name {
            Ok("description") => {
                description = field.text().await.map_err(|e| UploadError(e.to_string()))?;
            }
            Ok("image") => {
                let content_type = content_type?;
                if content_type != "image/jpeg" && content_type != "image/png" {
                    // println!("{:?}", content_type);
                    return Err(UploadError(
                        "Invalid image type. Only JPG and PNG are allowed.".into(),
                    ));
                }

                image_file = Some(
                    field
                        .bytes()
                        .await
                        .map_err(|e| UploadError(e.to_string()))?,
                );
            }
            _ => {}
        }
    }

    if description.is_empty() {
        return Err(UploadError("Missing description".into()));
    }

    let id = Uuid::new_v4();
    let compressed_id = Uuid::new_v4();
    match image_file {
        Some(img) => {
            let size = img.len();
            let result = state
                .s3_client
                .put_object()
                .bucket("media-service")
                .key(id)
                .body(img.into())
                .set_content_type(Some("application/jpeg".to_string()))
                .send()
                .await;

            match result {
                Ok(_output) => {
                    let _ = state
                        .rabbitmq_client
                        .send_message(MediaUploadedMessage {
                            id: id.to_string(),
                            compressed_id: compressed_id.to_string(),
                        })
                        .await;
                    Ok(Json(format!(
                        "Uploaded image with description: {} and Size {:?} kb",
                        description,
                        size / 1024
                    )))
                }
                Err(_err) => {
                    return Err(UploadError("Error uploading image".into()));
                }
            }
        }
        None => {
            return Err(UploadError("Missing image file.".into()));
        }
    }
}
