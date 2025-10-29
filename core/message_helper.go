package core

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// TuyaHeader represents the header of a Tuya message.
type TuyaHeader struct {
	Prefix     uint32
	Seqno      uint32
	Cmd        uint32
	Length     uint32
	TotalLength uint32
}

// MessagePayload represents a payload to be sent in a Tuya message.
type MessagePayload struct {
	Cmd     int
	Payload []byte
}

// TuyaMessage represents a Tuya message.
type TuyaMessage struct {
	Seqno    uint32
	Cmd      uint32
	Retcode  uint32
	Payload  []byte
	Crc      uint32
	CrcGood  bool
	Prefix   uint32
	IV       []byte
}

// PackMessage packs a TuyaMessage into bytes for protocol 3.3.
func PackMessage(msg TuyaMessage, key []byte) ([]byte, error) {
	var buffer bytes.Buffer

	// Encrypt payload for 3.3
	cipher := NewAESCipher(key)
	encryptedPayload, err := cipher.Encrypt(msg.Payload, false)
	if err != nil {
		return nil, err
	}

	// 3.3 header
	header := append([]byte(PROTOCOL_VERSION_BYTES_33), PROTOCOL_3x_HEADER...)
	buffer.Write(header)
	buffer.Write(encryptedPayload)

	// Recalculate payload length for header
	payloadLen := len(buffer.Bytes()) + 8 // 4 for CRC, 4 for suffix

	finalBuffer := new(bytes.Buffer)
	binary.Write(finalBuffer, binary.BigEndian, uint32(PREFIX_VALUE))
	binary.Write(finalBuffer, binary.BigEndian, msg.Seqno)
	binary.Write(finalBuffer, binary.BigEndian, msg.Cmd)
	binary.Write(finalBuffer, binary.BigEndian, uint32(payloadLen))

	finalBuffer.Write(buffer.Bytes())

	// CRC
	crc := crc32.ChecksumIEEE(finalBuffer.Bytes())
	binary.Write(finalBuffer, binary.BigEndian, crc)

	// Suffix
	binary.Write(finalBuffer, binary.BigEndian, uint32(SUFFIX_VALUE))

	return finalBuffer.Bytes(), nil
}

// UnpackMessage unpacks bytes into a TuyaMessage for protocol 3.3.
func UnpackMessage(data []byte, key []byte) (*TuyaMessage, error) {
	reader := bytes.NewReader(data)

	var header TuyaHeader
	binary.Read(reader, binary.BigEndian, &header.Prefix)
	binary.Read(reader, binary.BigEndian, &header.Seqno)
	binary.Read(reader, binary.BigEndian, &header.Cmd)
	binary.Read(reader, binary.BigEndian, &header.Length)

	if header.Prefix != PREFIX_VALUE {
		return nil, fmt.Errorf("invalid prefix")
	}

	payload := make([]byte, header.Length-8)
	reader.Read(payload)

	var crc uint32
	binary.Read(reader, binary.BigEndian, &crc)

	var suffix uint32
	binary.Read(reader, binary.BigEndian, &suffix)

	if suffix != SUFFIX_VALUE {
		return nil, fmt.Errorf("invalid suffix")
	}

	calculatedCrc := crc32.ChecksumIEEE(data[:len(data)-8])
	crcGood := calculatedCrc == crc

	// Decrypt payload for 3.3
	var decryptedPayload []byte
	if len(payload) > 0 {
		// Remove 3.3 header
		if bytes.HasPrefix(payload, []byte(PROTOCOL_VERSION_BYTES_33)) {
			payload = payload[len(PROTOCOL_VERSION_BYTES_33)+len(PROTOCOL_3x_HEADER):]
			cipher := NewAESCipher(key)
			var err error
			decryptedPayload, err = cipher.Decrypt(payload, false)
			if err != nil {
				// Try to decode as JSON if decryption fails
				return &TuyaMessage{
					Seqno:   header.Seqno,
					Cmd:     header.Cmd,
					Payload: payload, // return raw payload
					Crc:     crc,
					CrcGood: crcGood,
					Prefix:  header.Prefix,
				}, nil
			}
		} else {
			decryptedPayload = payload
		}
	}

	return &TuyaMessage{
		Seqno:   header.Seqno,
		Cmd:     header.Cmd,
		Payload: decryptedPayload,
		Crc:     crc,
		CrcGood: crcGood,
		Prefix:  header.Prefix,
	}, nil
}

// PackPlaintext55AA packs a TuyaMessage with a plaintext payload into a 55AA frame.
func PackPlaintext55AA(msg TuyaMessage) ([]byte, error) {
	payloadLen := len(msg.Payload) + 8 // 4 for CRC, 4 for suffix

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, uint32(PREFIX_VALUE))
	binary.Write(buffer, binary.BigEndian, msg.Seqno)
	binary.Write(buffer, binary.BigEndian, msg.Cmd)
	binary.Write(buffer, binary.BigEndian, uint32(payloadLen))
	buffer.Write(msg.Payload)

	// CRC
	crc := crc32.ChecksumIEEE(buffer.Bytes())
	binary.Write(buffer, binary.BigEndian, crc)

	// Suffix
	binary.Write(buffer, binary.BigEndian, uint32(SUFFIX_VALUE))

	return buffer.Bytes(), nil
}

// UnpackPlaintext55AA unpacks a 55AA frame with a plaintext payload.
func UnpackPlaintext55AA(data []byte) (*TuyaMessage, error) {
	reader := bytes.NewReader(data)

	var header TuyaHeader
	binary.Read(reader, binary.BigEndian, &header.Prefix)
	if header.Prefix != PREFIX_VALUE {
		// Not a 55AA message, try 6699
		if header.Prefix == PREFIX_6699_VALUE {
			// This is a 6699 message, should be handled elsewhere
			return nil, fmt.Errorf("unexpected 6699 prefix in plaintext unpacker")
		}
		return nil, fmt.Errorf("invalid 55AA prefix: %x", header.Prefix)
	}

	binary.Read(reader, binary.BigEndian, &header.Seqno)
	binary.Read(reader, binary.BigEndian, &header.Cmd)
	binary.Read(reader, binary.BigEndian, &header.Length)

	if len(data) < int(header.Length)+8 {
		return nil, fmt.Errorf("message too short")
	}

	payload := make([]byte, header.Length-8)
	reader.Read(payload)

	var crc uint32
	binary.Read(reader, binary.BigEndian, &crc)

	var suffix uint32
	binary.Read(reader, binary.BigEndian, &suffix)
	if suffix != SUFFIX_VALUE {
		return nil, fmt.Errorf("invalid 55AA suffix")
	}

	calculatedCrc := crc32.ChecksumIEEE(data[:len(data)-8])
	crcGood := calculatedCrc == crc

	return &TuyaMessage{
		Seqno:   header.Seqno,
		Cmd:     header.Cmd,
		Payload: payload,
		Crc:     crc,
		CrcGood: crcGood,
		Prefix:  header.Prefix,
	}, nil
}

// PackMessage6699 packs a TuyaMessage into bytes for protocol 3.5.
func PackMessage6699(msg TuyaMessage, sessionKey []byte) ([]byte, error) {
	// Generate random 12-byte IV
	iv := make([]byte, 12)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}

	// Pre-calculate length: IV (12) + payload + tag (16)
	length := uint32(12 + len(msg.Payload) + 16)

	// Construct AAD: reserved(2) + seq(4) + cmd(4) + length(4)
	aadBuffer := new(bytes.Buffer)
	binary.Write(aadBuffer, binary.BigEndian, uint16(0))
	binary.Write(aadBuffer, binary.BigEndian, msg.Seqno)
	binary.Write(aadBuffer, binary.BigEndian, msg.Cmd)
	binary.Write(aadBuffer, binary.BigEndian, length)
	aad := aadBuffer.Bytes()

	// Encrypt payload
	ciphertext, tag, err := GCMEncrypt(sessionKey, iv, msg.Payload, aad)
	if err != nil {
		return nil, err
	}

	// Assemble final packet
	finalBuffer := new(bytes.Buffer)
	binary.Write(finalBuffer, binary.BigEndian, uint32(PREFIX_6699_VALUE))
	finalBuffer.Write(aad)
	finalBuffer.Write(iv)
	finalBuffer.Write(ciphertext)
	finalBuffer.Write(tag)
	binary.Write(finalBuffer, binary.BigEndian, uint32(SUFFIX_6699_VALUE))

	return finalBuffer.Bytes(), nil
}

// UnpackMessage6699 unpacks bytes into a TuyaMessage for protocol 3.5.
func UnpackMessage6699(data []byte, sessionKey []byte) (*TuyaMessage, error) {
	reader := bytes.NewReader(data)

	var prefix uint32
	binary.Read(reader, binary.BigEndian, &prefix)
	if prefix != PREFIX_6699_VALUE {
		return nil, fmt.Errorf("invalid 6699 prefix")
	}

	// Read AAD fields
	var reserved uint16
	var seqno, cmd, length uint32
	binary.Read(reader, binary.BigEndian, &reserved)
	binary.Read(reader, binary.BigEndian, &seqno)
	binary.Read(reader, binary.BigEndian, &cmd)
	binary.Read(reader, binary.BigEndian, &length)

	// Construct AAD for decryption
	aadBuffer := new(bytes.Buffer)
	binary.Write(aadBuffer, binary.BigEndian, reserved)
	binary.Write(aadBuffer, binary.BigEndian, seqno)
	binary.Write(aadBuffer, binary.BigEndian, cmd)
	binary.Write(aadBuffer, binary.BigEndian, length)
	aad := aadBuffer.Bytes()

	iv := make([]byte, 12)
	reader.Read(iv)

	ciphertextAndTag := make([]byte, length-12)
	reader.Read(ciphertextAndTag)
	ciphertext := ciphertextAndTag[:len(ciphertextAndTag)-16]
	tag := ciphertextAndTag[len(ciphertextAndTag)-16:]

	var suffix uint32
	binary.Read(reader, binary.BigEndian, &suffix)
	if suffix != SUFFIX_6699_VALUE {
		return nil, fmt.Errorf("invalid 6699 suffix")
	}

	plaintext, err := GCMDecrypt(sessionKey, iv, ciphertext, tag, aad)
	if err != nil {
		return nil, fmt.Errorf("GCM decryption failed: %w", err)
	}

	// Handle return code
	var retcode uint32
	var payload []byte
	if len(plaintext) >= 4 {
		retcode = binary.BigEndian.Uint32(plaintext[:4])
		payload = plaintext[4:]
	} else {
		payload = plaintext
	}

	return &TuyaMessage{
		Seqno:   seqno,
		Cmd:     cmd,
		Retcode: retcode,
		Payload: payload,
		Prefix:  prefix,
		IV:      iv,
	}, nil
}
