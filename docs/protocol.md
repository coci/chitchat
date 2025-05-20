# 📡 ChitChat Binary Protocol Specification(v1)

## Overview

The **ChitChat Binary Protocol** is a compact, efficient, binary-based protocol used for communication between clients and the ChitChat server. It is designed to serve as the foundational layer of messaging and authentication in the ChitChat application. This protocol facilitates structured request–response interactions using a minimalistic message format that supports extensibility and efficient parsing.

All communication is carried out over a persistent TCP connection. Messages are sent in binary format and must conform to the structure described below.

---

## 🧱 Message Format

- **Header:**
    - `opcode`: 1 byte
    - `length`: 4 bytes (uint32, big-endian) — length of the body
    - `checksum`: 4 bytes — checksum for integrity verification

- **Body:**  
  Variable-length payload specific to each message type.

---

## Opcode Descriptions

### 🆕 Opcode: `0x01` — create user

**Direction:** Client → Server  
**Purpose:** Sends a request to create a new user account in the system.

**Body Structure:**  
The username and password of the user should be sent as a comma-separated string:
```
alice,123456
```
**Expected Response:**
- **Success:**
```
header opcode : 0x02
body : userid
```
- **Failure:**
```
header opcode : 0xFF
body : opcode of request,error message
```
for example :
```
header opcode : 0xFF
body : 0x01,User already exists
```