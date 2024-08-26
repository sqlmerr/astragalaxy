mod commands;

use lib_core::{
    create_mongodb_client,
    mongodb::Database,
    repositories::{location::MongoLocationRepository, user::MongoUserRepository},
    schemas::user::UserSchema,
    services::{location::LocationService, user::UserService},
};
use serenity::prelude::*;
use tracing::info;

type Error = Box<dyn std::error::Error + Send + Sync>;
type Context<'a> = poise::Context<'a, Data, Error>;
type ApplicationContext<'a> = poise::ApplicationContext<'a, Data, Error>;

pub struct Data {
    pub user_service: UserService<MongoUserRepository>,
    pub location_service: LocationService<MongoLocationRepository>,
}

pub struct InvocationData {
    pub user: UserSchema,
}

fn create_data(database: Database) -> Data {
    let user_repository = MongoUserRepository::new(database.collection("users"));
    let user_service = UserService::new(user_repository);

    let location_repository = MongoLocationRepository::new(database.collection("locations"));
    let location_service = LocationService::new(location_repository);

    Data {
        user_service,
        location_service,
    }
}

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok().unwrap();

    tracing_subscriber::fmt().init();

    let client = create_mongodb_client(std::env::var("MONGODB_URI").unwrap()).await;
    let database = client.database("astragalaxy");

    let options = poise::FrameworkOptions {
        command_check: Some(|ctx: Context| {
            Box::pin(async move {
                let user_id = ctx.author().id.get() as i64;

                let location = &ctx
                    .data()
                    .location_service
                    .find_one_location_by_code("space_station".to_string())
                    .await
                    .unwrap();

                let user_service = &ctx.data().user_service;
                let user = match user_service.find_one_user_by_discord_id(user_id).await {
                    Ok(u) => u,
                    Err(_) => {
                        let u = user_service
                            .register_from_discord(user_id, ctx.author().name.clone(), location._id)
                            .await;
                        match u {
                            Ok(_) => {}
                            Err(_) => return Ok(false),
                        }

                        u.unwrap()
                    }
                };

                ctx.set_invocation_data(InvocationData { user }).await;

                Ok(true)
            })
        }),
        commands: vec![commands::ping(), commands::location()],
        ..Default::default()
    };

    let framework = poise::Framework::builder()
        .setup(move |ctx, ready, framework| {
            Box::pin(async move {
                info!("Logged in as {}", ready.user.name);
                poise::builtins::register_globally(ctx, &framework.options().commands).await?;
                Ok(create_data(database))
            })
        })
        .options(options)
        .build();

    let token = std::env::var("DISCORD_BOT_TOKEN").expect("Expected a token in the environment");

    let intents = GatewayIntents::GUILD_MESSAGES
        | GatewayIntents::DIRECT_MESSAGES
        | GatewayIntents::MESSAGE_CONTENT;

    let mut client = Client::builder(&token, intents)
        .framework(framework)
        .await
        .expect("Err creating client");

    if let Err(why) = client.start().await {
        println!("Client error: {why:?}");
    }
}
