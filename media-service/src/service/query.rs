use crate::entity::{user_media, user_media::Entity as UserMedia};
use sea_orm::*;

pub struct Query;

impl Query {
    pub async fn find_post_by_media_id(
        db: &DbConn,
        media_id: &str,
    ) -> Result<Option<user_media::Model>, DbErr> {
        let user_media = UserMedia::find()
            .filter(user_media::Column::MediaId.contains(media_id))
            .one(db)
            .await?;
        Ok(user_media)
    }

    /// If ok, returns (post models, num pages).
    pub async fn find_posts_in_page(
        db: &DbConn,
        page: u64,
        posts_per_page: u64,
    ) -> Result<(Vec<user_media::Model>, u64), DbErr> {
        // Setup paginator
        let paginator = UserMedia::find()
            .order_by_asc(user_media::Column::Id)
            .paginate(db, posts_per_page);
        let num_pages = paginator.num_pages().await?;

        // Fetch paginated posts
        paginator.fetch_page(page - 1).await.map(|p| (p, num_pages))
    }
}
