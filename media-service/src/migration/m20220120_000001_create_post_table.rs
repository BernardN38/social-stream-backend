use sea_orm_migration::{
    prelude::*,
    schema::{integer, pk_auto, string},
};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager
            .create_table(
                Table::create()
                    .table(UserMedia::Table)
                    .if_not_exists()
                    .col(pk_auto(UserMedia::Id))
                    .col(integer(UserMedia::UserId))
                    .col(string(UserMedia::MediaId))
                    .col(string(UserMedia::MediaCompressedId))
                    .col(string(UserMedia::Status))
                    .to_owned(),
            )
            .await
    }

    async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager
            .drop_table(Table::drop().table(UserMedia::Table).to_owned())
            .await
    }
}

#[derive(DeriveIden)]
enum UserMedia {
    Table,
    Id,
    UserId,
    MediaId,
    MediaCompressedId,
    Status,
}
