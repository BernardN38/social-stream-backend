use amqprs::channel::{BasicConsumeArguments, QueueBindArguments, QueueDeclareArguments};
use aws_config::BehaviorVersion;
use aws_sdk_s3::config::{Credentials, Region};
mod rabbitmq_client;
use log::LevelFilter;
use rabbitmq_client::client::{CustomConsumer, RabbitmqClient, RabbitmqConfig};
use simplelog::SimpleLogger;
use std::env;
use tokio::signal;

#[tokio::main]
async fn main() {
    SimpleLogger::init(LevelFilter::Debug, simplelog::Config::default()).unwrap();

    let key_id = env::var("minioID").unwrap();
    let secret_key = env::var("minioAccessKey").unwrap();
    let cred = Credentials::new(key_id, secret_key, None, None, "loaded-from-custom-env");
    let s3_config = aws_sdk_s3::config::Builder::new()
        // .endpoint_resolver(ep)
        .endpoint_url("http://minio:9000")
        .behavior_version(BehaviorVersion::v2024_03_28())
        .credentials_provider(cred)
        .region(Region::new("us-west-1"))
        .force_path_style(true) // apply bucketname as path param instead of pre-domain
        .build();

    let client = aws_sdk_s3::Client::from_conf(s3_config);
    let rabbitmq_client = RabbitmqClient::new(
        RabbitmqConfig {
            host: "rabbitmq".to_owned(),
            port: 5672,
            username: "guest".to_owned(),
            password: "guest".to_owned(),
            exchange: "media_events".to_owned(),
        },
        client.clone(),
    )
    .await
    .unwrap();

    println!("consuming started");
    let channel = rabbitmq_client.connection.open_channel(None).await.unwrap();
    let (queue_name, _, _) = channel
        .queue_declare(
            QueueDeclareArguments::default()
                .queue("media_compresion_service".to_owned())
                .durable(true)
                .finish(),
        )
        .await
        .unwrap()
        .unwrap();
    channel
        .queue_bind(QueueBindArguments::new(
            &queue_name,
            &rabbitmq_client.config.exchange,
            "media.uploaded",
        ))
        .await
        .unwrap();
    // start consumer with given name
    let args = BasicConsumeArguments::new(&queue_name, "example_basic_pub_sub");

    let consumer = CustomConsumer {
        s3_client: client.clone(),
    };
    channel.basic_consume(consumer, args).await.unwrap();
    println!("awaitng shutdown signal");
    shutdown_signal().await;
    println!("shutdown signal received");
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
