pub mod error;
pub mod schemas;

use base64::prelude::BASE64_STANDARD;
use base64::Engine;
use ed25519_dalek::{Signature, Verifier, VerifyingKey};
pub use error::*;
use hmac::{Hmac, Mac};
use once_cell::sync::Lazy;
use rand::Rng;
pub use schemas::*;
use sha2::{Digest, Sha256};
use std::collections::HashMap;
use std::time::{Duration, SystemTime, UNIX_EPOCH};
use tokio::time::timeout;
use tonlib::cell::BagOfCells;
use tonlib::client::TonClient;
use tonlib::config::{MAINNET_CONFIG, TESTNET_CONFIG};
use tonlib::contract::{TonContractFactory, TonContractInterface};
use tonlib::wallet::{
    WalletDataHighloadV2R2, WalletDataV1V2, WalletDataV3, WalletDataV4, WalletVersion,
};

pub const PAYLOAD_TTL: u64 = 3600; // 1 hour
pub const PROOF_TTL: u64 = 3600; // 1 hour
const KNOWN_HASHES: Lazy<HashMap<[u8; 32], WalletVersion>> = Lazy::new(|| {
    let mut known_hashes = HashMap::new();
    let all_versions = [
        WalletVersion::V1R1,
        WalletVersion::V1R2,
        WalletVersion::V1R3,
        WalletVersion::V2R1,
        WalletVersion::V2R2,
        WalletVersion::V3R1,
        WalletVersion::V3R2,
        WalletVersion::V4R1,
        WalletVersion::V4R2,
        WalletVersion::HighloadV1R1,
        WalletVersion::HighloadV1R2,
        WalletVersion::HighloadV2,
        WalletVersion::HighloadV2R1,
        WalletVersion::HighloadV2R2,
    ];
    all_versions.into_iter().for_each(|v| {
        let hash: [u8; 32] = v
            .code()
            .unwrap()
            .cell_hash()
            .try_into()
            .expect("all hashes [u8; 32], right?");
        known_hashes.insert(hash, v);
    });
    known_hashes
});

pub async fn generate_ton_proof_payload(secret: String) -> Result<GenerateTonProofPayload, Error> {
    let mut payload: [u8; 48] = [0; 48];
    let mut rng = rand::thread_rng();
    rng.fill(&mut payload[8..8]);

    if let Ok(n) = SystemTime::now().duration_since(UNIX_EPOCH) {
        let expire = n.as_secs() + PAYLOAD_TTL;
        let expire_be = expire.to_be_bytes();
        payload[8..16].copy_from_slice(&expire_be);
    } else {
        return Err(Error("SystemTime before UNIX_EPOCH".to_string()).into());
    }

    let mut mac = Hmac::<Sha256>::new_from_slice(secret.as_bytes())
        .map_err(|_| Error("Error creating hmac".to_string()))?;
    mac.update(&payload[0..16]);
    let signature = mac.finalize().into_bytes();
    payload[16..48].copy_from_slice(&signature);

    let hex = hex::encode(&payload[0..32]);

    Ok(GenerateTonProofPayload { payload: hex })
}

pub async fn check_ton_proof(
    payload: CheckProofPayload,
    secret: String,
    domain: String,
) -> Result<String, Error> {
    let data = hex::decode(payload.proof.payload.clone())
        .map_err(|_| Error("Error decoding payload".to_string()))?;

    if data.len() != 32 {
        return Err(Error("Invalid payload length".to_string()));
    }
    let mut mac = Hmac::<Sha256>::new_from_slice(secret.as_bytes())
        .map_err(|_| Error("Error creating hmac".to_string()))?;
    mac.update(&data[..16]);
    let signature_bytes: [u8; 32] = mac.finalize().into_bytes().into();

    let signature_valid = data
        .iter()
        .skip(16)
        .zip(signature_bytes.iter().take(16))
        .all(|(a, b)| a == b);
    if !signature_valid {
        return Err(Error("invalid payload signature".to_string()));
    }

    let now = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map_err(|_| Error("SystemTime before UNIX_EPOCH".to_string()))?
        .as_secs();

    // check payload expiration
    let expire_b: [u8; 8] = data[8..16].try_into().expect("already checked length");
    let expire_d = u64::from_be_bytes(expire_b);

    if now > expire_d {
        return Err(Error("payload expired".to_string()));
    }

    if now > payload.proof.timestamp + PROOF_TTL {
        return Err(Error("ton proof has been expired".to_string()));
    }

    if payload.proof.domain.value != domain {
        return Err(Error(format!("wrong domain. got {}, expected {}", payload.proof.domain.value, domain)));
    }

    if payload.proof.domain.length_bytes != payload.proof.domain.value.len() as u64 {
        return Err(Error(
            "domain length is not equal to domain value length".to_string(),
        ));
    }

    let ton_proof_prefix = "ton-proof-item-v2/";
    let mut msg: Vec<u8> = Vec::new();
    msg.extend_from_slice(ton_proof_prefix.as_bytes());
    msg.extend_from_slice(&payload.address.workchain.to_be_bytes());
    msg.extend_from_slice(&payload.address.hash_part);
    msg.extend_from_slice(&(payload.proof.domain.length_bytes as u32).to_le_bytes());
    msg.extend_from_slice(payload.proof.domain.value.as_bytes());
    msg.extend_from_slice(&payload.proof.timestamp.to_le_bytes());
    msg.extend_from_slice(payload.proof.payload.as_bytes());

    let mut hasher = Sha256::new();
    hasher.update(msg);
    let msg_hash = hasher.finalize();

    let mut full_msg: Vec<u8> = vec![0xff, 0xff];
    let ton_connect_prefix = "ton-connect";
    full_msg.extend_from_slice(ton_connect_prefix.as_bytes());
    full_msg.extend_from_slice(&msg_hash);

    let mut hasher = Sha256::new();
    hasher.update(full_msg);
    let full_msg_hash = hasher.finalize();

    let client = match payload.network {
        TonNetwork::Mainnet => TonClient::builder()
            .with_config(MAINNET_CONFIG)
            .build()
            .await
            .map_err(|_| Error("Error creating ton client".to_string()))?,
        TonNetwork::Testnet => TonClient::builder()
            .with_config(TESTNET_CONFIG)
            .build()
            .await
            .map_err(|_| Error("Error creating ton client".to_string()))?,
    };

    let contract_factory = TonContractFactory::builder(&client)
        .build()
        .await
        .map_err(|_| Error("Server error".to_string()))?;
    let wallet_contract = contract_factory.get_contract(&payload.address);

    let pubkey_bytes = match timeout(
        Duration::from_secs(10),
        wallet_contract.run_get_method("get_public_key", &vec![]),
    )
    .await
    {
        Ok(Ok(r)) => {
            let pubkey_n = r.stack[0]
                .get_biguint()
                .map_err(|_| Error("server error".to_string()))?;
            let pubkey_b: [u8; 32] = pubkey_n.to_bytes_be().try_into().map_err(|_| {
                Error("failed to extract 32 bits long public from the wallet contract".to_string())
            })?;
            pubkey_b
        }
        Err(_) | Ok(Err(_)) => {
            let bytes = BASE64_STANDARD
                .decode(&payload.proof.state_init)
                .map_err(|_| Error("Decode error".to_string()))?;
            let boc = BagOfCells::parse(&bytes).map_err(|e| Error(e.to_string()))?;
            let hash: [u8; 32] = boc
                .single_root()
                .map_err(|_| Error("ton cell error".to_string()))?
                .cell_hash()
                .try_into()
                .map_err(|_| Error("invalid state_init length".to_string()))?;

            if hash != payload.address.hash_part {
                return Err(Error("wrong address in state_init".to_string()));
            }

            let root = boc.single_root().expect("checked above");
            let code = root.reference(0).map_err(|e| Error(e.to_string()))?;
            let data = root
                .reference(1)
                .map_err(|e| Error(e.to_string()))?
                .as_ref()
                .clone();

            let code_hash: [u8; 32] = code
                .cell_hash()
                .try_into()
                .map_err(|_| Error("invalid code of wallet".to_string()))?;
            let version = KNOWN_HASHES
                .get(&code_hash)
                .ok_or(Error("not known wallet version".to_string()))?
                .clone();

            let pubkey_b = match version {
                WalletVersion::V1R1
                | WalletVersion::V1R2
                | WalletVersion::V1R3
                | WalletVersion::V2R1
                | WalletVersion::V2R2 => {
                    let data = WalletDataV1V2::try_from(data).map_err(|e| Error(e.to_string()))?;
                    data.public_key
                }
                WalletVersion::V3R1 | WalletVersion::V3R2 => {
                    let data = WalletDataV3::try_from(data).map_err(|e| Error(e.to_string()))?;
                    data.public_key
                }
                WalletVersion::V4R1 | WalletVersion::V4R2 => {
                    let data = WalletDataV4::try_from(data)
                        .map_err(|_| Error("server error".to_string()))?;
                    data.public_key
                }
                WalletVersion::HighloadV2R2 => {
                    let data =
                        WalletDataHighloadV2R2::try_from(data).map_err(|e| Error(e.to_string()))?;
                    data.public_key
                }
                _ => {
                    return Err(Error("can't process given wallet version".to_string()));
                }
            };

            pubkey_b
        }
    };

    let pubkey =
        VerifyingKey::from_bytes(&pubkey_bytes).map_err(|_| Error("invalid pubkey".to_string()))?;
    let signature_bytes: [u8; 64] = BASE64_STANDARD
        .decode(&payload.proof.signature)
        .map_err(|e| Error(e.to_string()))?
        .try_into()
        .map_err(|_| Error("expected 64 bit long signature".to_string()))?;
    let signature = Signature::from_bytes(&signature_bytes);
    pubkey
        .verify(&full_msg_hash, &signature)
        .map_err(|e| Error(e.to_string()))?;

    Ok(payload.address.to_base64_std())
}
