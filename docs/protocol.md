# ChitChat Binary Message Protocol (CBMP)

**Author**: Sorush Safari 
**Version**: 1.0  
**Date**: 2025-05-26

---

## 1. Introduction

This document defines the **Simple Binary Message Protocol (SBMP)**, a lightweight protocol for encoding and decoding binary messages consisting of a fixed-size header and variable-length body. It is designed for efficient transmission of structured messages over byte streams such as TCP connections.

---

## 2. Message Structure

Each message consists of a **Header** followed by a **Body**.

```
+--------+--------------------+---------------------+
| Opcode |   Length (uint32)  |      Body[...]      |
| (1B)   |   (4B, BigEndian)  | (Length bytes long) |
+--------+--------------------+---------------------+
```

### 2.1 Header

- **Opcode (1 byte)**: Identifies the message type.
- **Length (4 bytes, unsigned int, big-endian)**: Specifies the number of bytes in the body.

### 2.2 Body

- **Body (N bytes)**: A comma-separated string payload. The size is determined by the `Length` field in the header.

  > **Note:** The body MUST be a UTF-8 encoded string with comma-separated values.

for example :
```text
john,doe
```
---

## 3. Encoding

To encode a message:

1. Set the `Opcode` field to the appropriate message type identifier.
2. Set the `Length` field to the length of the `Body`.
3. Allocate a buffer of size `5 + Length`.
4. Write the header (`Opcode` + `Length`) into the first 5 bytes.
5. Copy the `Body` into the buffer immediately following the header.

---

## 4. Decoding

To decode a message:

1. Read 5 bytes from the stream to obtain the header.
2. Parse `Opcode` and `Length`.
3. Read `Length` bytes from the stream to obtain the body.

---

## 5. Validation

After decoding, a message is considered **valid** if the following conditions are met:

- `Length` > 0
- The size of the body equals the `Length` specified in the header.

If validation fails, the message MUST be rejected.

---

## 6. Error Handling

- If reading the header fails (e.g., connection closed), return a "read header failed" error.
- If reading the body fails (e.g., incomplete data), return a "read body failed" error.
- If `Length` does not match the size of `Body`, return a "body length mismatch" error.

---

## 7. Example

### Example Message

| Field  | Value        |
| ------ | ------------ |
| Opcode | `0x01`       |
| Length | `0x00000009` |
| Body   | `john,doe`   |

### Binary Encoding (Hex):

```
01 00 00 00 09 6A 6F 68 6E 2C 64 6F 65
```

---

## 8. Extensions

This protocol can be extended by:

- Defining semantics for `Opcode` values (e.g., `0x01 = PING`, `0x02 = MESSAGE`, etc.)
- Including optional fields or checksums in the body if needed.
- Adding compression/encryption outside the protocol scope.

---

## 9. Security Considerations

- Always validate `Length` before allocating memory to avoid resource exhaustion.
- Consider limiting the maximum `Length` (e.g., max 1MB) to protect against large payload attacks.