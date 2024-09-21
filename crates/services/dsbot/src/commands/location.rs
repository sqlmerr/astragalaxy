use lib_core::schemas::user::UserSchema;
use lib_utils::snake_case_to_normal;

use poise::CreateReply;
use serenity::all::{
    colours, ButtonStyle, CreateActionRow, CreateButton, CreateEmbed, CreateInteractionResponse,
    CreateInteractionResponseMessage,
};

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
    let players_count = ctx
        .data()
        .user_service
        .get_users_count_by_location(user.location_id)
        .await;

    let reply = {
        let components = vec![CreateActionRow::Buttons(vec![CreateButton::new("button")
            .label("Button")
            .style(ButtonStyle::Primary)])];

        let title;
        if location.multiplayer {
            title = format!("{location_name} - {players_count} players");
        } else {
            title = location_name
        }
        CreateReply::default()
            .embed(
                CreateEmbed::new()
                    .title(title)
                    .color(colours::branding::FUCHSIA),
            )
            .components(components)
    };

    ctx.send(reply).await?;

    while let Some(interaction) =
        poise::serenity_prelude::ComponentInteractionCollector::new(ctx.serenity_context())
            .timeout(std::time::Duration::from_secs(120))
            .filter(move |interaction| interaction.data.custom_id == "button")
            .await
    {
        match interaction
            .create_response(
                &ctx,
                CreateInteractionResponse::UpdateMessage(
                    CreateInteractionResponseMessage::default()
                        .add_embed(CreateEmbed::new().field("sus", "asas", false))
                        .button(CreateButton::new_link("https://google.com").label("sus")),
                ),
            )
            .await
        {
            Err(_) => return Ok(()),
            Ok(_) => {}
        };
    }

    Ok(())
}
