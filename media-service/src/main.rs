use aws_config::BehaviorVersion;
use aws_sdk_s3::config::{Credentials, Region};
use axum::{
    extract::DefaultBodyLimit,
    routing::{get, post},
    Router,
};
use handlers::handlers::{check_health, upload_media};
use rabbitmq_client::client::{RabbitmqClient, RabbitmqConfig};
use std::env;
use tokio::signal;
mod handlers;
mod jwt;
mod rabbitmq_client;
use tower_http::limit::RequestBodyLimitLayer;

#[derive(Clone, Debug)]
pub struct AppState {
    s3_client: aws_sdk_s3::Client,
    rabbitmq_client: RabbitmqClient,
}

#[tokio::main]
async fn main() {
    // let app_config = get_env_config();
    let app_config = match get_env_config() {
        Ok(config) => config,
        Err(e) => {
            println!("{:?}", e);
            return;
        }
    };

    let rabbitmq_client = RabbitmqClient::new(RabbitmqConfig {
        host: "rabbitmq".to_owned(),
        port: 5672,
        username: "guest".to_owned(),
        password: "guest".to_owned(),
        exchange: "media_events".to_owned(),
    })
    .await
    .unwrap();

    let cred = Credentials::new(
        app_config.minio_id,
        app_config.minio_secret_key,
        None,
        None,
        "loaded-from-custom-env",
    );
    let s3_config = aws_sdk_s3::config::Builder::new()
        // .endpoint_resolver(ep)
        .endpoint_url("http://minio:9000")
        .behavior_version(BehaviorVersion::v2024_03_28())
        .credentials_provider(cred)
        .region(Region::new("us-west-1"))
        .force_path_style(true) // apply bucketname as path param instead of pre-domain
        .build();

    let client = aws_sdk_s3::Client::from_conf(s3_config);

    let state = AppState {
        s3_client: client,
        rabbitmq_client: rabbitmq_client,
    };
    // build our application with a single route
    let app = Router::new()
        .route("/api/v1/media/health", get(check_health))
        .route("/api/v1/media/upload", post(upload_media)) // Adding the middleware
        .layer(DefaultBodyLimit::disable()) // Disable default limit to manage it manually
        .layer(RequestBodyLimitLayer::new(50 * 1024 * 1024)) // 50 MB limit/ Handle errors (see below)
        // .layer(TraceLayer::new_for_http()) // Optional: Log requests
        .with_state(state);

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    println!("{}", "listening on post 8080".to_string());
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

async fn shutdown_signal() {
    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install signal handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }
}
#[derive(Debug)]
struct AppConfig {
    minio_id: String,
    minio_secret_key: String,
}
#[derive(Debug)]
pub struct ConfigError(String);
fn get_env_config() -> Result<AppConfig, ConfigError> {
    let mut app_config = AppConfig {
        minio_id: "".to_string(),
        minio_secret_key: "".to_string(),
    };
    app_config.minio_id = env::var("minioID").map_err(|e| return ConfigError(e.to_string()))?;
    app_config.minio_secret_key =
        env::var("minioAccessKey").map_err(|e| return ConfigError(e.to_string()))?;
    return Ok(app_config);
}
