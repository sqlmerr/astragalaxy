use axum::{
    extract::{Query, State},
    http::StatusCode,
    middleware,
    routing::{get, post},
    Extension, Json, Router,
};
use lib_auth::errors::AuthError;
use lib_core::{
    errors::CoreError,
    schemas::{
        spaceship::CreateSpaceshipSchema,
        user::{CreateUserSchema, UpdateUserSchema, UserSchema},
    },
};
use lib_utils::parse_token;

use crate::{
    errors::Result,
    middlewares::auth::{auth_middleware, protection_middleware},
    schemas::auth::{AuthBody, AuthPayload, GetTokenSchema, UserTokenSchema},
    state::ApplicationState,
};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = middleware::from_fn_with_state(state.clone(), auth_middleware);
    let protection_middleware = middleware::from_fn_with_state(state, protection_middleware);

    Router::new()
        .route("/register", post(register))
        .route(
            "/register/telegram",
            post(register_from_telegram).layer(protection_middleware.clone()),
        )
        .route(
            "/token/sudo",
            get(get_user_token_telegram_sudo).layer(protection_middleware),
        )
        .route("/login", post(login))
        // .route("/discord", post(discord_callback))
        .route("/me", get(profile).layer(auth_middleware))
}

async fn register(
    State(state): State<ApplicationState>,
    Json(user): Json<CreateUserSchema>,
) -> Result<(StatusCode, Json<UserSchema>)> {
    let location = state
        .location_service
        .find_one_location_by_code("space_station".to_string())
        .await?;

    let system = state
        .system_service
        .find_one_system_by_name(String::from("initial_system"))
        .await?;
    let user = state
        .user_service
        .register(user, location._id, system._id)
        .await?;
    let spaceship = state
        .spaceship_service
        .create_spaceship(CreateSpaceshipSchema {
            name: "initial".to_string(),
            user_id: user._id,
            location_id: location._id,
            system_id: system._id,
        })
        .await?;

    state
        .user_service
        .update_user(
            user._id,
            UpdateUserSchema {
                username: None,
                spaceship_id: Some(spaceship._id),
            },
        )
        .await?;

    let user = state.user_service.find_one_user(user._id).await?;

    Ok((StatusCode::CREATED, Json(user)))
}

async fn login(
    State(state): State<ApplicationState>,
    Json(payload): Json<AuthPayload>,
) -> Result<Json<AuthBody>> {
    let (telegram_id, token) = match parse_token(payload.token) {
        Some(t) => t,
        None => return Err(CoreError::AuthError(AuthError::InvalidToken).into()),
    };

    let jwt_token: String = state.user_service.login(telegram_id, token).await?;

    Ok(Json(AuthBody::new(jwt_token)))
}

async fn profile(Extension(user): Extension<UserSchema>) -> Json<UserSchema> {
    Json(user)
}

// async fn discord_callback(
//     State(state): State<ApplicationState>,
//     Json(payload): Json<DiscordAuthPayload>,
// ) -> Result<Json<AuthBody>> {
//     let response: DiscordAuthBody = request(
//         "https://discord.com/api/v10/oauth2/token".to_string(),
//         "POST".parse().unwrap(),
//     )
//     .form(&json!(
//         {
//             "grant_type": "authorization_code",
//             "code": payload.code,
//             "redirect_uri": format!("{}/auth/discord", state.config.domain)
//         }
//     ))
//     .basic_auth(
//         state.config.discord_client_id,
//         Some(state.config.discord_client_secret),
//     )
//     .send()
//     .await
//     .map_err(|_| CoreError::ServerError)?
//     .error_for_status()
//     .map_err(|e| {
//         println!("{:?}", e);
//         CoreError::ServerError
//     })?
//     .json()
//     .await
//     .map_err(|_| CoreError::ServerError)?;

//     println!("{:?}", serde_json::to_string_pretty(&response).unwrap());

//     // return Err(CoreError::ServerError.into());

//     let discord_user: Value = request(
//         "https://discord.com/api/v10/users/@me".to_string(),
//         "GET".parse().unwrap(),
//     )
//     .bearer_auth(response.access_token)
//     .send()
//     .await
//     .map_err(|_| CoreError::ServerError)?
//     .error_for_status()
//     .map_err(|e| {
//         println!("{:?}", e);
//         CoreError::ServerError
//     })?
//     .json()
//     .await
//     .map_err(|_| CoreError::ServerError)?;

//     println!("{}", serde_json::to_string_pretty(&discord_user).unwrap());

//     let user = state
//         .user_service
//         .find_one_user_by_discord_id(discord_user["id"].as_str().unwrap().parse::<i64>().unwrap())
//         .await;
//     println!("{:?}", user);

//     let user = match user {
//         Err(_) => {
//             let location = state
//                 .location_service
//                 .find_one_location_by_code("space_station".to_string())
//                 .await?;
//             let system = state
//                 .system_service
//                 .find_one_system_by_name(String::from("initial_system"))
//                 .await?;

//             state
//                 .user_service
//                 .register_from_discord(
//                     discord_user["id"].as_str().unwrap().parse::<i64>().unwrap(),
//                     discord_user["username"].as_str().unwrap().to_string(),
//                     location._id,
//                     system._id,
//                 )
//                 .await?
//         }
//         Ok(user) => user,
//     };
//     let token = create_token(&Claims::new(user.username))
//         .map_err(|_| CoreError::AuthError(AuthError::TokenCreation))?;

//     Ok(Json(AuthBody::new(token)))
// }

async fn register_from_telegram(
    State(state): State<ApplicationState>,
    Json(user): Json<CreateUserSchema>,
) -> Result<(StatusCode, Json<UserSchema>)> {
    let location = state
        .location_service
        .find_one_location_by_code("space_station".to_string())
        .await?;

    let system = state
        .system_service
        .find_one_system_by_name(String::from("initial_system"))
        .await?;

    let user = state
        .user_service
        .register(user, location._id, system._id)
        .await?;

    let spaceship = state
        .spaceship_service
        .create_spaceship(CreateSpaceshipSchema {
            name: "initial".to_string(),
            user_id: user._id,
            location_id: location._id,
            system_id: system._id,
        })
        .await?;

    state
        .user_service
        .update_user(
            user._id,
            UpdateUserSchema {
                username: None,
                spaceship_id: Some(spaceship._id),
            },
        )
        .await?;

    let user = state.user_service.find_one_user(user._id).await?;

    Ok((StatusCode::CREATED, Json(user)))
}

async fn get_user_token_telegram_sudo(
    State(state): State<ApplicationState>,
    Query(payload): Query<GetTokenSchema>,
) -> Result<Json<UserTokenSchema>> {
    let user = state
        .user_service
        .find_one_raw_user_by_telegram_id(payload.telegram_id)
        .await?;

    let token = user.token;

    Ok(Json(UserTokenSchema { token }))
}
