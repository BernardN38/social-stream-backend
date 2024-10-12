use amqprs::{
    callbacks::{DefaultChannelCallback, DefaultConnectionCallback},
    channel::{
        BasicAckArguments, BasicConsumeArguments, BasicNackArguments, BasicPublishArguments,
        BasicRejectArguments, Channel, ExchangeDeclareArguments, QueueBindArguments,
        QueueDeclareArguments,
    },
    connection::{self, Connection, OpenConnectionArguments},
    consumer::{AsyncConsumer, DefaultConsumer},
    BasicProperties, Deliver,
};
use core::fmt;
use std::{
    f32::consts::E,
    fmt::{Debug, Formatter},
};

use crate::rabbitmq_client::models::MediaUploadedMessage;
#[derive(Clone)]
pub struct RabbitmqClient {
    connection: Connection,
    channel: Channel,
    config: RabbitmqConfig,
}

impl Debug for RabbitmqClient {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        f.debug_struct("RabbitmqClient")
            .field("connection", &self.channel.to_string())
            .field("channel", &self.channel.to_string())
            .field("config", &self.config)
            .finish()
    }
}
#[derive(Clone, Debug)]
pub struct RabbitmqConfig {
    pub host: String,
    pub port: u16,
    pub username: String,
    pub password: String,

    pub exchange: String,
}
#[derive(Debug)]
pub struct RabbitError(String);
impl RabbitmqClient {
    pub async fn new(config: RabbitmqConfig) -> Result<Self, RabbitError> {
        let connection = Connection::open(&OpenConnectionArguments::new(
            &config.host,
            config.port,
            &config.username,
            &config.password,
        ))
        .await;
        let connection = match connection {
            Ok(c) => c,
            Err(e) => return Err(RabbitError(e.to_string())),
        };
        // open a channel on the connection
        let channel = connection.open_channel(None).await.unwrap();

        let channel_result = channel
            .exchange_declare(
                ExchangeDeclareArguments::default()
                    .exchange("media_events".to_owned())
                    .exchange_type("topic".to_owned())
                    .finish(),
            )
            .await;
        match channel_result {
            Ok(c) => {}
            Err(e) => return Err(RabbitError(e.to_string())),
        }
        return Ok(Self {
            connection: connection,
            channel: channel,
            config: config,
        });
    }
    pub async fn send_message(&self, message: impl Into<Vec<u8>>) -> Result<(), RabbitError> {
        // create arguments for basic_publish
        let args = BasicPublishArguments::new(&self.config.exchange, "media.uploaded");

        self.channel
            .basic_publish(BasicProperties::default(), message.into(), args)
            .await
            .map_err(|e| return RabbitError(e.to_string()))
    }
}

pub struct CustomConsumer {}
impl AsyncConsumer for CustomConsumer {
    fn consume<'life0, 'life1, 'async_trait>(
        &'life0 mut self,
        channel: &'life1 Channel,
        deliver: Deliver,
        basic_properties: BasicProperties,
        content: Vec<u8>,
    ) -> ::core::pin::Pin<
        Box<dyn ::core::future::Future<Output = ()> + ::core::marker::Send + 'async_trait>,
    >
    where
        'life0: 'async_trait,
        'life1: 'async_trait,
        Self: 'async_trait,
    {
        Box::pin(async move {
            // Process the message content
            let msg = serde_json::from_slice(&content);
            let msg = match msg {
                Ok(m) => m,
                Err(e) => {
                    let _ = channel
                        .basic_reject(BasicRejectArguments::new(deliver.delivery_tag(), false))
                        .await;
                    println!("{:?}", e);
                    return;
                }
            };
            // msg.await;
            // println!("Received message: {:?}", data);

            // Acknowledge the message (depending on your use case)
            if let Err(err) = channel
                .basic_ack(BasicAckArguments::new(deliver.delivery_tag(), false))
                .await
            {
                eprintln!("Failed to ack message: {:?}", err);
            }
        })
    }
}
