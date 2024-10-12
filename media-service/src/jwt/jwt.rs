use axum::{
    async_trait,
    body::Body,
    extract::{FromRequestParts, Request, State},
    http::{
        request::{self, Parts},
        StatusCode,
    },
    middleware::Next,
    response::{IntoResponse, Response},
    Json,
};
use chrono::Utc;
use jsonwebtoken::{decode, DecodingKey, Validation};
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::fmt::Display;

static KEYS: Lazy<Keys> = Lazy::new(|| {
    let secret = std::env::var("jwtSecret").expect("JWT_SECRET must be set");
    Keys::new(secret.as_bytes())
});

impl Display for Claims {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "User id: {}\nExpires: {}", self.user_id, self.exp)
    }
}

#[async_trait]
impl<S> FromRequestParts<S> for Claims
where
    S: Send + Sync, // Service state must live at least as long as 'a
{
    type Rejection = AuthError;
    async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
        // Extract the cookies from the request
        let cookie_header = parts.headers.get(axum::http::header::COOKIE);
        let cookie_header = match cookie_header {
            Some(header) => header,
            None => return Err(AuthError::MissingToken),
        };
        // Parse the cookies
        let cookies = cookie_header
            .to_str()
            .map_err(|_| AuthError::InvalidToken)?;

        if let Some(jwt) = extract_jwt(cookies) {
            let validation = Validation::new(jsonwebtoken::Algorithm::HS512);

            // Decode the user data
            let token_data = decode::<Claims>(jwt, &KEYS.decoding, &validation);

            match token_data {
                Ok(t) => {
                    // println!(
                    //     "token expires at: {:?}, current time: {}",
                    //     &t.claims.exp,
                    //     Utc::now().timestamp()
                    // );
                    Ok(t.claims)
                }
                Err(e) => {
                    println!("{}", e);
                    Err(AuthError::InvalidToken)
                }
            }
        } else {
            println!("JWT not found");
            return Err(AuthError::MissingToken);
        }
    }
}
pub fn extract_jwt(cookie_str: &str) -> Option<&str> {
    // Find the start of the "jwt=" substring
    if let Some(start_index) = cookie_str.find("jwt=") {
        // Slice the string starting from the character after "jwt="
        let start = start_index + "jwt=".len();

        // Find the next semicolon after "jwt=", indicating the end of the token
        // Use Option::map_or to handle the case where there's no semicolon
        let end = cookie_str[start..]
            .find(';')
            .map_or(cookie_str.len() - start, |semicolon_index| semicolon_index)
            + start;

        // Return the slice of the JWT value
        Some(&cookie_str[start..end])
    } else {
        None
    }
}

impl IntoResponse for AuthError {
    fn into_response(self) -> Response {
        println!("making response, {:?}", self);
        let (status, error_message) = match self {
            AuthError::InvalidToken => (StatusCode::BAD_REQUEST, "Invalid token"),
            AuthError::MissingToken => (StatusCode::BAD_REQUEST, "Missing token"),
        };
        let body = Json(json!({
            "error": error_message
        }));

        let response = (status, body).into_response();
        response
    }
}

struct Keys {
    decoding: DecodingKey,
}

impl Keys {
    fn new(secret: &[u8]) -> Self {
        Self {
            decoding: DecodingKey::from_secret(secret),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub user_id: i32,
    pub exp: usize,
}

#[derive(Debug)]
pub enum AuthError {
    InvalidToken,
    MissingToken,
}

// pub async fn auth_guard(
//     // parts: Parts,
//     req: Request,
//     next: Next,
// ) -> Result<Response<Body>, AuthError> {
//     // println!("running middleware, {:?}", parts.headers);
//     // // Extract the cookies from the request
//     // req.headers().get(key)
//     let cookie_header = req.headers().get(axum::http::header::COOKIE);
//     let cookie_header = match cookie_header {
//         Some(header) => header,
//         None => return Err(AuthError::MissingToken),
//     };

//     // println!("looking for cookie");
//     // Parse the cookies
//     let cookies = cookie_header
//         .to_str()
//         .map_err(|_| AuthError::InvalidToken)?;

//     if let Some(jwt) = extract_jwt(cookies) {
//         let validation = Validation::new(jsonwebtoken::Algorithm::HS512);

//         // Decode the user data
//         let token_data = decode::<Claims>(jwt, &KEYS.decoding, &validation);

//         match token_data {
//             Ok(t) => {
//                 // println!(
//                 //     "token expires at: {:?}, current time: {}",
//                 //     &t.claims.exp,
//                 //     Utc::now().timestamp()
//                 // );
//                 // Ok(t.claims);
//             }
//             Err(e) => {
//                 println!("{}", e);
//                 return Err(AuthError::InvalidToken);
//             }
//         }
//     } else {
//         println!("JWT not found");
//         return Err(AuthError::MissingToken);
//     }
//     let result = next.run(req).await;
//     return Ok(result);
// }
