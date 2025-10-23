# 🗣️ GOSSIP Protocol Specification

```
Version: 1.0
Purpose: Binary protocol for secure, low-latency communication in chat systems.
Transport Layer: TCP over TLS 1.3
Endian: Big-endian (network byte order)
Authentication: Username + Password → Session (persistent device-bound tokens)
Integrity: HMAC-SHA256 truncated to 16 bytes
Confidentiality: TLS 1.3 encryption
```

## 1. Overview

GOSSIP is a binary application-level protocol designed for chat and messaging systems.
It provides:
 - Secure session establishment (with TLS 1.3)
 - Long-lived device sessions
 - Frame-based binary messages with fixed headers
 - Strong integrity via per-session HMAC
 - Replay and ordering protection using sequence numbers
 - Simple extensibility via message typesq

#### 1.3 Design Philosophy

GOSSIP is minimal. Every byte has a purpose. It uses two magic bytes for framing synchronization and integrity, followed by a fixed header and a variable payload. The body supports rich binary without requiring schema negotiation.


## 2. Transport and Security
• All communication must occur over TLS 1.3.

• Server presents a valid certificate.

• ALPN identifier: "gossip/1".

• Recommended ciphers:

- TLS_AES_128_GCM_SHA256

- TLS_CHACHA20_POLY1305_SHA256

• Optional: mutual TLS for administrative or service clients.

No application-level encryption beyond HMAC integrity is required.
HMAC protects against replay, duplication, or session confusion.



## 3. Frame Structure

Every GOSSIP message is a Frame consisting of a fixed header and a variable-length payload.

#### 3.1 Short Header (Pre-Auth)

Used before authentication — e.g., login, hello, token renew.
```text

+-------------+--------------------------------------------------------------+
|    Header   |  Magic     |  version    |  MsgType  |  Stream ID |  Length  |
|   12 bytes  |  0xNU 0xLL |  1 byte     |  1 byte   |  4 bytes   |  2 bytes |             
|+------------+------------+------------+------------+------------+----------+
|     Body    |                      (Length bytes)                          |
+-------------+--------------------------------------------------------------+
```

• Header size: 10 bytes

• Payload: raw text

• Integrity: Transport-only (TLS).

#### 3.2 Long Header (Post-Auth / Secure)

Used after successful authentication. Adds session identity, sequence tracking, and HMAC tag.

```text
+-------------+------------------------------------------------------------------------------------------------+
|    Header   |  Magic     |  version    |  MsgType  |  Stream ID |  Length  | session_id | seq     | tag      |
|   12 bytes  |  0xNU 0xLL |  1 byte     |  1 byte   |  4 bytes   |  2 bytes |  4 bytes   | 8 bytes | 16 bytes |
|+------------+------------+------------+------------+------------+----------+------------+---------+----------+
|     Body    |                                        (Length bytes)                                          | 
+-------------+------------------------------------------------------------------------------------------------+
```

• Header size: 38 bytes

• Payload: raw text

• Integrity: Application-level HMAC

## 4. HMAC Integrity Verification
```text
tag = first16bytes( HMAC_SHA256(session_key, header_without_tag || payload) )
```

• session_key: 32-byte random secret assigned at login

• Verification:
Lookup session_key by session_id.

• Compute HMAC

• Compare constant-time

• Drop frame if invalid or seq ≤ last_seen.

## 5. Sequence Numbers
• 64-bit unsigned integer (seq)

• Increment by 1 for every frame per direction

• Server maintains seq_up (client→server) and seq_down (server→client)

• Drop or resync if sequence regresses or jumps too far

• Prevents replay or duplication

## 6. Session Model

| Concept | Description                                         |
|--------|-----------------------------------------------------|
|    user| Registered account                                  |
|    device_id| Random UUID per installation                        |
|    session_id| 32-bit server-side handle linking user + device     |
| session_key| 256-bit random secret for HMAC; never leaves server |
| refresh_token| Long-lived (e.g., 365 days), rotated each use.      |
|   access_token| Short-lived (e.g., 15 min), optional                |
|   seq| Monotonic counter per direction                     |


## 7. Authentication Lifecycle
1.	Login (AUTH_LOGIN_REQ → AUTH_LOGIN_OK):
 - Validates credentials over TLS
 - Server issues session_id and session_key
2.	Normal Operation:
 -	Use Long Header for all frames
 -	Include session_id, seq, tag
3.	Token Renewal:
 - Send TOKEN_RENEW_REQ with refresh_token (Short Header, still under TLS)
 - Server replies TOKEN_RENEW_OK with rotated tokens
4.	Logout:
 - Client sends LOGOUT_REQ
 - Server deletes session; returns LOGOUT_OK