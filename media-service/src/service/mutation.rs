use crate::entity::{user_media, user_media::Entity as UserMedia};
use sea_orm::*;

pub struct Mutation;

impl Mutation {
    pub async fn create_post(
        db: &DbConn,
        form_data: user_media::Model,
    ) -> Result<user_media::ActiveModel, DbErr> {
        user_media::ActiveModel {
            user_id: Set(form_data.user_id),
            media_id: Set(form_data.media_id),
            media_compressed_id: Set(form_data.media_compressed_id),
            status: Set(form_data.status),
            ..Default::default()
        }
        .save(db)
        .await
    }

    pub async fn update_user_media_by_id(
        db: &DbConn,
        media_id: &str,
        status: String,
    ) -> Result<user_media::Model, DbErr> {
        let user_media = UserMedia::find()
            .filter(user_media::Column::MediaId.contains(media_id))
            .one(db)
            .await?;
        match user_media {
            Some(user_media) => {
                user_media::ActiveModel {
                    id: Set(user_media.id),
                    user_id: Set(user_media.user_id),
                    media_id: Set(user_media.media_id),
                    media_compressed_id: Set(user_media.media_compressed_id),
                    status: Set(status),
                }
                .update(db)
                .await
            }
            None => Err(sea_orm::DbErr::ConvertFromU64("none user media")),
        }
    }

    // pub async fn delete_post(db: &DbConn, id: i32) -> Result<DeleteResult, DbErr> {
    //     let post: user_media::ActiveModel = Post::find_by_id(id)
    //         .one(db)
    //         .await?
    //         .ok_or(DbErr::Custom("Cannot find post.".to_owned()))
    //         .map(Into::into)?;

    //     post.delete(db).await
    // }

    // pub async fn delete_all_posts(db: &DbConn) -> Result<DeleteResult, DbErr> {
    //     Post::delete_many().exec(db).await
    // }
}
