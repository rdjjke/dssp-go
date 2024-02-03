# Message layout

Each message starts with a header followed by a content.

Header layout can be found in a table format: [message_layout.xls](message_layout.xls).

Notes:
- To decide, which endianness is used, the receiver checks the first message before processing it. If the third bit of 
  the last byte in CTLWORD = 1, it's Little Endian. If the third bit of the first byte = 1, it's Big Endian.
- Only message content can be compressed / encrypted / authenticated.
- Content can be sent only in the second message of the client / server.

## Handshake

Notes:
- If during a handshake, any content option is set, it only means that this option can be used in the stream, but it's
not obligatory.
- Nonce is the concatenation of SEQ + NONCESUF fields.
- Session key (used to authenticate messages after the handshake) and encryption key are calculated as HKDF(pre-shared 
  key, client nonce, server nonce). First 32 bytes are used for authentication, last 16 - for encryption.
- Crypto sizes: 
  - Pre-shared key is of 32 bytes.
  - Nonces are of 12 bytes (4-byte TIME + 8-byte random NONCESUF).
  - Session key is of 32 bytes (used for HMAC-SHA256).
  - Encryption key is of 16 bytes (used for AES-128).
- For congestion and flow control, sending/receiving rates are used instead of window sizes. This approach is more 
  suitable for long fat networks.
