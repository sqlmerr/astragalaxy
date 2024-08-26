use lib_core::schemas::user::UserSchema;
use lib_utils::snake_case_to_normal;

use poise::CreateReply;
use serenity::all::{colours, CreateEmbed};

use crate::{Context, Error, InvocationData};

#[poise::command(slash_command, track_edits)]
pub async fn location(ctx: Context<'_>) -> Result<(), Error> {
    let user: UserSchema = ctx
        .invocation_data::<InvocationData>()
        .await
        .as_deref()
        .unwrap()
        .user
        .clone();

    let location = ctx
        .data()
        .location_service
        .find_one_location(user.location_id)
        .await?;
    let location_name = snake_case_to_normal(location.code.as_str());

    ctx.send(
        CreateReply::default()
            .embed(
                CreateEmbed::new()
                    .title(format!("You are in {}", location_name))
                    .color(colours::branding::FUCHSIA),
            )
            .ephemeral(true),
    )
    .await?;
    Ok(())
}
