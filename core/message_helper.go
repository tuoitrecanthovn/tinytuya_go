package core

import (
	"bytes"
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
