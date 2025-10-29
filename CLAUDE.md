# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of the TinyTuya protocol for local communication with Tuya smart home devices. It supports protocol versions 3.1 through 3.5, with automatic version detection and device discovery capabilities.

## Build and Run Commands

```bash
# Build the application
go build -o tinytuya_go

# Run with default configuration (tests all protocol versions)
go run main.go

# Clean build cache if compilation issues occur
go clean -cache && go mod tidy

# Build for production (no CGO required)
CGO_ENABLED=0 go build -o tinytuya_go
```

## Architecture Overview

### Core Components

**XenonDevice** (`core/xenon_device.go`): The main device implementation handling:
- Session key negotiation (v3.4+)
- Protocol version detection
- TCP connection management
- Message send/receive operations

**Message Helpers** (`core/message_helper.go`): Protocol-specific message formatting:
- 55AA frame format (v3.1-3.4): `prefix | seqno | cmd | length | payload | CRC32 | suffix`
- 6699 frame format (v3.5): `prefix | reserved | seq | cmd | length | IV | ciphertext | tag | suffix`

**Crypto Layer** (`core/crypto_helper.go`): Security implementations:
- AES-ECB encryption (v3.1-3.4)
- AES-GCM encryption (v3.5)
- Session key derivation (v3.4+)
- HMAC-SHA256 integrity (v3.4)

**Device Types**: Automatic detection and handling:
- Default: Standard device communication
- device22: For devices requiring explicit DPS mapping (22-character gwId)
- Gateway support for sub-devices

### Protocol Version Support

The implementation automatically tests protocol versions in order: 3.3 → 3.4 → 3.5

- **v3.1**: Basic AES-ECB, MD5 validation, monitoring without key
- **v3.3**: AES-ECB encryption, version headers, CRC32 integrity
- **v3.4**: Session key negotiation (3-way handshake), HMAC-SHA256 integrity
- **v3.5**: AES-GCM per-packet encryption, global sequencing

### Session Key Negotiation (v3.4+)

Critical for successful communication with newer devices:

1. **START (0x03)**: Client sends 16-byte nonce
2. **RESP (0x04)**: Device responds with encrypted device_nonce + HMAC(client_nonce)
3. **FINISH (0x05)**: Client sends HMAC(device_nonce)
4. **Key Derivation**: XOR nonces → encrypt → extract session key slice

## Key Implementation Details

### Network Configuration
- **TCP**: Port 6668 for device communication
- **UDP**: Ports 6666/6667 for discovery, 7000 for solicited discovery (v3.5)
- **Timeouts**: 10-second connection timeout recommended
- **Retry Logic**: Built-in retry mechanism for failed connections

### Device Discovery
- Automatic IP detection when Address = "Auto" or "0.0.0.0"
- UDP broadcast discovery on local network
- Solicited discovery for v3.5 devices

### DPS (Data Point System)
Device state is represented as JSON with DPS keys:
```json
{
  "dps": {
    "1": true,      // Switch state
    "2": 24.5,      // Temperature sensor
    "3": "auto"     // Mode setting
  }
}
```

### Error Handling Patterns
- "data unvalid" response triggers device22 mode automatically
- EOF errors typically indicate protocol version mismatch or network issues
- Session negotiation failures fall back to earlier protocol versions

## Debugging Common Issues

### Connection Failures
1. Verify device IP and local key (16 characters for v3.4+)
2. Check network connectivity and firewall rules
3. Test with Python TinyTuya as reference implementation
4. Try protocol version 3.3 if v3.5 fails

### Session Key Negotiation Issues
- Ensure local key is exactly 16 characters for v3.4+
- Check if device actually supports claimed protocol version
- Monitor for EOF errors during handshake (indicates version mismatch)

### Device Type Detection
- Devices with 22-character gwId require device22 mode
- Library auto-detects and switches to device22 on "data unvalid" response
- Manual device type override available via NewXenonDevice parameters

## Contrib Modules

Device-specific implementations in `contrib/` package:
- Climate devices, thermostats, doorbells
- Specialized DPS mappings and control logic
- Example: `climate_device.go` for air conditioner control

## Testing Strategy

The main.go includes automated compatibility testing:
1. Tests protocol versions sequentially
2. Performs basic control operations (switch toggle)
3. Validates device state changes
4. Includes device22 mode testing for compatible devices

No unit test framework currently exists - integration testing via main.go.