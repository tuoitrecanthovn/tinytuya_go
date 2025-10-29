package main

import (
	"fmt"
	"log"
	"time"

	"tinytuya_go/core"
)

func main() {
	fmt.Println("tinytuya_go v3.5 example")

	const (
		DeviceID    = "ebaf701f96b67923b5vna9" // REPLACE WITH YOUR v3.5 DEVICE ID
		IPAddress   = "192.168.1.119"          // REPLACE WITH YOUR DEVICE IP
		LocalKey    = "Cq|@r(278+lDpzah"         // REPLACE WITH YOUR LOCAL KEY
		ProtocolVer = 3.5
	)

	d, err := core.NewXenonDevice(
		DeviceID,
		IPAddress,
		LocalKey,
		"default", // devType
		5*time.Second,
		ProtocolVer,
		false, // persist
		"",    // cid
		nil,   // parent
	)

	if err != nil {
		log.Fatalf("Failed to create device: %v", err)
	}

	log.Println("Getting device status...")
	data, err := d.Status()
	if err != nil {
		log.Fatalf("Failed to get device status: %v", err)
	}

	log.Printf("Device status: %+v\n", data)

	// Example: Turn on a switch (assuming DPS '1' is the switch)
	log.Println("\nTurning device ON (DPS 1)...")
	_, err = d.SetValue("1", true)
	if err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	// Wait a moment
	time.Sleep(2 * time.Second)

	log.Println("\nGetting updated status...")
	data, err = d.Status()
	if err != nil {
		log.Fatalf("Failed to get device status: %v", err)
	}
	log.Printf("Device status: %+v\n", data)

	log.Println("\nTurning device OFF (DPS 1)...")
	_, err = d.SetValue("1", false)
	if err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	log.Println("\nExample finished.")
}
