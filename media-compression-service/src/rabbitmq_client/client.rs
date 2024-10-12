use amqprs::{
    channel::{BasicAckArguments, BasicPublishArguments, Channel, ExchangeDeclareArguments},
    connection::{Connection, OpenConnectionArguments},
    consumer::AsyncConsumer,
    BasicProperties, Deliver,
};
use aws_sdk_s3::Client;
use image::{imageops::FilterType, GenericImageView};
use log::{error, info};
use std::{
    fmt::Debug,
    time::{Duration, Instant},
};
use tokio::time::sleep;
use turbojpeg::{PixelFormat, Subsamp};

use crate::rabbitmq_client::models::{MediaCompressedMessage, MediaUploadedMessage};

#[derive(Clone)]
pub struct RabbitmqClient {
    s3_client: aws_sdk_s3::Client,
    pub channel: Channel,
    pub connection: Connection,
    pub config: RabbitmqConfig,
}

#[derive(Clone, Debug)]
pub struct RabbitmqConfig {
    pub host: String,
    pub port: u16,
    pub username: String,
    pub password: String,
    pub exchange: String,
}

impl RabbitmqClient {
    pub async fn new(
        config: RabbitmqConfig,
        s3_client: aws_sdk_s3::Client,
    ) -> Result<Self, String> {
        let connection = Connection::open(&OpenConnectionArguments::new(
            &config.host,
            config.port,
            &config.username,
            &config.password,
        ))
        .await
        .map_err(|e| format!("Error opening connection: {}", e))?;

        let channel = connection
            .open_channel(None)
            .await
            .map_err(|e| format!("Error opening channel: {}", e))?;

        channel
            .exchange_declare(
                ExchangeDeclareArguments::default()
                    .exchange("media_events".to_owned())
                    .exchange_type("topic".to_owned())
                    .finish(),
            )
            .await
            .map_err(|e| format!("Error declaring exchange: {}", e))?;

        Ok(Self {
            connection,
            channel,
            config,
            s3_client,
        })
    }

    pub async fn send_message(&self, message: impl Into<Vec<u8>>) -> Result<(), String> {
        let args = BasicPublishArguments::new(&self.config.exchange, "media.uploaded");

        self.channel
            .basic_publish(BasicProperties::default(), message.into(), args)
            .await
            .map_err(|e| format!("Failed to publish message: {}", e))
    }
}

pub struct CustomConsumer {
    pub s3_client: Client,
}

// impl AsyncConsumer for CustomConsumer {
//     fn consume<'life0, 'life1, 'async_trait>(
//         &'life0 mut self,
//         channel: &'life1 Channel,
//         deliver: Deliver,
//         basic_properties: BasicProperties,
//         content: Vec<u8>,
//     ) -> ::core::pin::Pin<
//         Box<dyn ::core::future::Future<Output = ()> + ::core::marker::Send + 'async_trait>,
//     >
//     where
//         'life0: 'async_trait,
//         'life1: 'async_trait,
//         Self: 'async_trait,
//     {
//         Box::pin(async move {
//             let data: MediaUploadedMessage = match serde_json::from_slice(&content) {
//                 Ok(data) => data,
//                 Err(e) => {
//                     error!("Failed to parse message: {}", e);
//                     return;
//                 }
//             };
//             info!("Message received: {:?}", data);
//             sleep(Duration::from_secs(2)).await;
//             let start = Instant::now();
//             match self
//                 .s3_client
//                 .get_object()
//                 .bucket("media-service")
//                 .key(&data.id)
//                 .send()
//                 .await
//             {
//                 Ok(img) => {
//                     // info!("{:?}", img);
//                     match img.content_length() {
//                         Some(size) => {
//                             if size < 8 * 1024 * 1024 {
//                                 info!("image small enough no compression needed");
//                                 if let Err(err) = channel
//                                     .basic_ack(BasicAckArguments::new(
//                                         deliver.delivery_tag(),
//                                         false,
//                                     ))
//                                     .await
//                                 {
//                                     error!("Failed to acknowledge message: {}", err);
//                                 }
//                                 return;
//                             }
//                         }
//                         None => return,
//                     }
//                     let img_bytes = match img.body.collect().await {
//                         Ok(b) => b,
//                         Err(e) => {
//                             error!("Failed to collect image data: {}", e);
//                             return;
//                         }
//                     };

//                     let compressed_img = match compress_image(img_bytes.to_vec()) {
//                         Ok(img) => img,
//                         Err(e) => {
//                             error!("Failed to compress image: {:?}", e);
//                             return;
//                         }
//                     };

//                     match self
//                         .s3_client
//                         .put_object()
//                         .bucket("media-service")
//                         .key(&data.compressed_id)
//                         .body(compressed_img.into())
//                         .content_type("application/jpeg")
//                         .send()
//                         .await
//                     {
//                         Ok(_) => info!("Image successfully compressed and uploaded"),
//                         Err(e) => error!("Failed to upload compressed image: {}", e),
//                     }
//                 }
//                 Err(e) => error!("Failed to retrieve image from S3: {}", e),
//             }
//             let end = Instant::now();
//             info!("Image compression duration: {:?}", end - start);

//             if let Err(err) = channel
//                 .basic_ack(BasicAckArguments::new(deliver.delivery_tag(), false))
//                 .await
//             {
//                 error!("Failed to acknowledge message: {}", err);
//             }
//         })
//     }
// }

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
        let s3_client = self.s3_client.clone();
        let channel = channel.clone();
        Box::pin(async move {
            // Spawn a new task to process the message concurrently
            tokio::spawn(async move {
                // Process the message
                let data: MediaUploadedMessage = match serde_json::from_slice(&content) {
                    Ok(data) => data,
                    Err(e) => {
                        error!("Failed to parse message: {}", e);
                        return;
                    }
                };
                info!("Message received: {:?}", data);
                sleep(Duration::from_secs(2)).await;
                match s3_client
                    .get_object()
                    .bucket("media-service")
                    .key(&data.id)
                    .send()
                    .await
                {
                    Ok(img) => {
                        match img.content_length() {
                            Some(size) => {
                                if size < 8 * 1024 * 1024 {
                                    info!("image small enough no compression needed");
                                    if let Err(err) = channel
                                        .basic_ack(BasicAckArguments::new(
                                            deliver.delivery_tag(),
                                            false,
                                        ))
                                        .await
                                    {
                                        error!("Failed to acknowledge message: {}", err);
                                    }
                                    return;
                                }
                            }
                            None => return,
                        }
                        let img_bytes = match img.body.collect().await {
                            Ok(b) => b,
                            Err(e) => {
                                error!("Failed to collect image data: {}", e);
                                return;
                            }
                        };

                        let compressed_img = match compress_image(img_bytes.to_vec()) {
                            Ok(img) => img,
                            Err(e) => {
                                error!("Failed to compress image: {:?}", e);
                                return;
                            }
                        };

                        match s3_client
                            .put_object()
                            .bucket("media-service")
                            .key(&data.compressed_id)
                            .body(compressed_img.into())
                            .send()
                            .await
                        {
                            Ok(_) => {
                                info!("Image successfully compressed and uploaded");
                                let publish_result = channel
                                    .basic_publish(
                                        BasicProperties::default(),
                                        MediaCompressedMessage {
                                            id: data.id,
                                            compressed_id: data.compressed_id,
                                            status: "compressed".to_string(),
                                        }
                                        .into(),
                                        BasicPublishArguments::default()
                                            .exchange("media_events".to_string())
                                            .routing_key("media.compressed".to_string())
                                            .finish(),
                                    )
                                    .await;
                                match publish_result {
                                    Ok(sucess) => return,
                                    Err(e) => error!("{}", e.to_string()),
                                }
                            }
                            Err(e) => error!("Failed to upload compressed image: {}", e),
                        }
                    }
                    Err(e) => error!("Failed to retrieve image from S3: {}", e),
                }

                // Acknowledge the message once processed
                if let Err(err) = channel
                    .basic_ack(BasicAckArguments::new(deliver.delivery_tag(), false))
                    .await
                {
                    error!("Failed to acknowledge message: {}", err);
                }
            });
        })
    }
}

#[derive(Debug)]
pub struct CompressionError(String);

fn compress_image(image_data: Vec<u8>) -> Result<Vec<u8>, CompressionError> {
    let image = image::load_from_memory(&image_data)
        .map_err(|e| CompressionError(format!("Failed to map image to dynamic image: {}", e)))?;

    let (width, height) = image.dimensions();
    let aspect_ratio = height as f32 / width as f32;
    let new_height = (1080 as f32 * aspect_ratio).round() as u32;

    let resized_img = image.resize(1080, new_height, FilterType::Lanczos3);
    let binding = resized_img.to_rgb8();
    let tj_image = turbojpeg::Image {
        pixels: binding.as_raw().as_slice(),
        width: resized_img.width() as usize,
        pitch: (resized_img.width() * PixelFormat::RGB.size() as u32) as usize,
        height: resized_img.height() as usize,
        format: PixelFormat::RGB,
    };
    turbojpeg::compress(tj_image, 80, Subsamp::Sub2x1)
        .map_err(|e| CompressionError(format!("Failed to compress image: {}", e)))
        .map(|compressed_image| compressed_image.to_vec())
}
