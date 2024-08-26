use std::net::SocketAddr;

use axum::{
    extract::{
        ws::{Message, WebSocket},
        ConnectInfo, Query, State, WebSocketUpgrade,
    },
    middleware,
    response::IntoResponse,
    routing::get,
    Extension, Router,
};
use lib_core::{errors::CoreError, schemas::user::UserSchema};
use serde::{Deserialize, Serialize};

use crate::{middlewares::auth::auth_middleware, state::ApplicationState};

pub(super) fn router(state: ApplicationState) -> Router<ApplicationState> {
    let auth_middleware = middleware::from_fn_with_state(state, auth_middleware);

    Router::new()
        .route("/ws", get(ws_handler))
        .route("/ws/", get(ws_handler))
        .layer(auth_middleware)
}

#[derive(Clone, Serialize, Deserialize)]
struct MoveData {
    x: i64,
    y: i64,
}

async fn ws_handler(
    Extension(user): Extension<UserSchema>,
    State(state): State<ApplicationState>,
    ws: WebSocketUpgrade,
    // ConnectInfo(addr): ConnectInfo<SocketAddr>,
) -> impl IntoResponse {
    println!("Websocket client connected.");

    ws.on_upgrade(move |socket| websocket(socket, state, user))
}

async fn websocket(
    mut socket: WebSocket,
    // who: SocketAddr,
    state: ApplicationState,
    user: UserSchema,
) {
    while let Some(msg) = socket.recv().await {
        let msg = match msg {
            Ok(msg) => msg,
            Err(err) => {
                println!("Websocket client disconnected: {err:?}");
                return;
            }
        };

        let message_text = match msg.to_text() {
            Err(_) => continue,
            Ok(text) => text,
        };

        let message = match serde_json::from_str::<serde_json::Value>(&message_text) {
            Ok(message) => message,
            Err(err) => {
                eprintln!("Error parsing message from client: {err:?}");
                continue;
            }
        };

        let response = match message.get("action") {
            Some(action) => match action.as_str() {
                Some("ping") => serde_json::json!({ "action": "ping", "data": "pong" }),
                Some("move") => {
                    let data: Result<MoveData, _> =
                        serde_json::from_value(message.get("data").unwrap().clone())
                            .map_err(|_| CoreError::ServerError);
                    let data = match data {
                        Ok(data) => data,
                        Err(_) => continue,
                    };

                    let x = data.x;
                    let y = data.y;

                    match state.user_service.move_user(user._id, x, y).await {
                        Ok(_) => serde_json::json!({ "action": "move", "success": true }),
                        Err(_) => serde_json::json!({ "action": "move", "success": false }),
                    }
                }
                _ => serde_json::json!({ "action": "unknown", "data": "unknown action" }),
            },
            None => serde_json::json!({ "action": "unknown", "data": "no action field" }),
        };

        if socket
            .send(Message::Text(serde_json::to_string(&response).unwrap()))
            .await
            .is_err()
        {
            println!("Websocket client disconnected");
            return;
        }
    }
}
