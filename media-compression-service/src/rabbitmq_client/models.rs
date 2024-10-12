use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, Debug)]
pub struct MediaUploadedMessage {
    pub id: String,
    pub compressed_id: String,
}

// impl Serialize for MediaUploadedMessage {
//     fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
//     where
//         S: serde::Serializer,
//     {
//         let mut state = serializer.serialize_struct("MediaUploadedMessage", 2)?;
//         state.serialize_field("id", &self.id.to_string())?;
//         state.serialize_field("compressed_id", &self.compressed_id.to_string())?;
//         state.end()
//     }
// }

impl Into<Vec<u8>> for MediaUploadedMessage {
    fn into(self) -> Vec<u8> {
        serde_json::to_vec(&self).unwrap()
    }
}
