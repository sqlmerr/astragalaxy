use poise::CreateReply;
use serenity::all::{colours, CreateEmbed};

use crate::{Context, Error};

#[poise::command(slash_command, track_edits)]
pub async fn ping(ctx: Context<'_>) -> Result<(), Error> {
    ctx.send(
        CreateReply::default()
            .embed(
                CreateEmbed::new()
                    .title("Pong!")
                    .color(colours::branding::RED)
                    .field("Random number", "2", false),
            )
            .ephemeral(true),
    )
    .await?;
    Ok(())
}
