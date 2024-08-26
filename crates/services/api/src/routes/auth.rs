use axum::{
    extract::State,
    http::StatusCode,
    middleware,
    routing::{get, post},
    Extension, Json, Router,
};
use lib_auth::{errors::AuthError, jwt::create_token, schemas::Claims};
use lib_core::{
    errors::CoreError,
    schemas::user::{CreateUserSchema, UserSchema},
};
use lib_utils::request;
use serde_json::{json, Value};

use crate::{
    errors::Result,
    middlewares::auth::auth_middleware,
    schemas::auth::{AuthBody, AuthPayload, DiscordAuthBody, DiscordAuthPayload},
    state::ApplicationState,
};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .route("/discord", post(discord_callback))
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
    let user = state.user_service.register(user, location._id).await?;

    Ok((StatusCode::CREATED, Json(user)))
}

async fn login(
    State(state): State<ApplicationState>,
    Json(payload): Json<AuthPayload>,
) -> Result<Json<AuthBody>> {
    let token = state
        .user_service
        .login(payload.username, payload.password)
        .await?;

    Ok(Json(AuthBody::new(token)))
}

async fn profile(Extension(user): Extension<UserSchema>) -> Json<UserSchema> {
    Json(user)
}

async fn discord_callback(
    State(state): State<ApplicationState>,
    Json(payload): Json<DiscordAuthPayload>,
) -> Result<Json<AuthBody>> {
    let response: DiscordAuthBody = request(
        "https://discord.com/api/v10/oauth2/token".to_string(),
        "POST".parse().unwrap(),
    )
    .form(&json!(
        {
            "grant_type": "authorization_code",
            "code": payload.code,
            "redirect_uri": format!("{}/auth/discord", state.config.domain)
        }
    ))
    .basic_auth(
        state.config.discord_client_id,
        Some(state.config.discord_client_secret),
    )
    .send()
    .await
    .map_err(|_| CoreError::ServerError)?
    .error_for_status()
    .map_err(|e| {
        println!("{:?}", e);
        CoreError::ServerError
    })?
    .json()
    .await
    .map_err(|_| CoreError::ServerError)?;

    println!("{:?}", serde_json::to_string_pretty(&response).unwrap());

    // return Err(CoreError::ServerError.into());

    let discord_user: Value = request(
        "https://discord.com/api/v10/users/@me".to_string(),
        "GET".parse().unwrap(),
    )
    .bearer_auth(response.access_token)
    .send()
    .await
    .map_err(|_| CoreError::ServerError)?
    .error_for_status()
    .map_err(|e| {
        println!("{:?}", e);
        CoreError::ServerError
    })?
    .json()
    .await
    .map_err(|_| CoreError::ServerError)?;

    println!("{}", serde_json::to_string_pretty(&discord_user).unwrap());

    let user = state
        .user_service
        .find_one_user_by_discord_id(discord_user["id"].as_str().unwrap().parse::<i64>().unwrap())
        .await;
    println!("{:?}", user);

    let user = match user {
        Err(_) => {
            let location = state
                .location_service
                .find_one_location_by_code("space_station".to_string())
                .await?;
            println!("hh");

            state
                .user_service
                .register_from_discord(
                    discord_user["id"].as_str().unwrap().parse::<i64>().unwrap(),
                    discord_user["username"].as_str().unwrap().to_string(),
                    location._id,
                )
                .await?
        }
        Ok(user) => user,
    };
    let token = create_token(&Claims::new(user.username))
        .map_err(|_| CoreError::AuthError(AuthError::TokenCreation))?;

    Ok(Json(AuthBody::new(token)))
}
