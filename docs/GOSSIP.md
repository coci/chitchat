# 🗣️ GOSSIP Protocol Specification

```
Version: 1.0
Transport: TCP
Endianness: Big-Endian
Document Type: Specification
```

## 1. Overview

#### 1.1 Purpose

GOSSIP is a lightweight, binary, stream-oriented protocol designed for efficient communication in distributed chat systems.
It provides a structured, extensible, and low-latency message exchange between clients, servers, and services such as gateways, routers, or message buses.

#### 1.2 Goals
•	🔹 Low overhead: compact binary format with fixed-length headers.
•	🔹 Extensible: flexible TLV (Type-Length-Value) payload encoding.
•	🔹 Reliable: built on TCP to ensure message delivery and order.
•	🔹 Simple to parse: predictable header layout and clear semantics.

#### 1.3 Design Philosophy

GOSSIP is minimal: every byte has a purpose. It uses two magic bytes for framing synchronization and integrity, followed by a fixed header and a variable payload. The TLV body supports rich structured data (such as JSON-like fields, binary attachments, or nested records) without requiring schema negotiation.


## 2. Frame Structure
```text
+------------------------------------------------------------------------------------+
|                                 Header (12 bytes)                                  |
+------------+------------+------------+------------+-----------------+--------------+
|  Magic     |  version    |  MsgType    |  Flags      |  Stream ID    |  Length     |
|  0xFU 0xCK |  1 byte     |  1 byte     |  2 bytes    |  4 bytes      |  2 bytes    |
+------------+------------+------------+------------+-----------------+--------------+
|                                 Payload (Length bytes)                             |
+------------------------------------------------------------------------------------+
```

Each message transmitted over GOSSIP is called a frame.
```text
┌──────────┬────────┬────────────┬─────────────────────────────────────────────────────────────────────────────────
│ Offset   │ Size   │ Field Name │ Description                         │ Type   │ Notes                            │
├──────────┼────────┼────────────┼─────────────────────────────────────├────────├──────────────────────────────────├ 
│ 0x00     │ 2 B    │ Magic      │ Frame marker for sync. (0x42 0x47)  │ uint16 │ Always same for protocol version │
│ 0x02     │ 1 B    │ Version    │ Protocol version                    │ uint8  │ Start with 0x01                  │
│ 0x03     │ 1 B    │ Msg Type   │ Opcode or message category          │ uint8  │ e.g., 0x01 = Chat, 0x02 = Join   │
│ 0x04     │ 2 B    │ Flags      │ Optional bits                       │ uint8  │ e.g., bit 0 = compressed         │
│ 0x05     │ 4 B    │ Stream ID  │ Conversation or channel identifier  │ uint32 │ Correlates request/response      │
│ 0x09     │ 2 B    │ Length     │ Payload size in bytes               │ uint16 │ Up to 65535                      │
│ 0x0B     │ N B    │ Payload    │ TLV structured payload              │ var    │ Decoded according to TLV rules   │
└──────────┴────────┴────────────┴─────────────────────────────────────────────────────────────────────────────────├
```
Total header size: 12 bytes
Payload: variable length (0–65,535 bytes)