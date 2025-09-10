package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GbxRemoteGo/mockserver"
)

func main() {
	// Create and start mock server
	mock := mockserver.New(mockserver.Config{
		Host: "127.0.0.1",
		Port: 5000,
	})

	// Set custom response for TriggerModeScriptEventArray
	mock.SetResponse("TriggerModeScriptEventArray", "Custom response: Pause activated")

	// Start the mock server
	if err := mock.Start(); err != nil {
		log.Fatal("Failed to start mock server:", err)
	}
	defer mock.Stop()

	fmt.Println("Mock server started on port", mock.Port())

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Create and connect client
	client := gbxclient.NewGbxClient("127.0.0.1", 5000, gbxclient.Options{})

	// Connect to the mock server
	if err := client.Connect(); err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Disconnect()

	// Authenticate
	if err := client.SetApiVersion("2023-04-24"); err != nil {
		log.Fatal("Failed to set API version:", err)
	}

	if err := client.EnableCallbacks(true); err != nil {
		log.Fatal("Failed to enable callbacks:", err)
	}

	if err := client.Authenticate("SuperAdmin", "SuperAdmin"); err != nil {
		log.Fatal("Failed to authenticate:", err)
	}

	fmt.Println("Client connected and authenticated successfully!")

	// Test some method calls
	version, err := client.GetVersion()
	if err != nil {
		log.Fatal("Failed to get version:", err)
	}
	fmt.Printf("Server version: %+v\n", version)

	status, err := client.GetStatus()
	if err != nil {
		log.Fatal("Failed to get status:", err)
	}
	fmt.Printf("Server status: %+v\n", status)

	// Test the custom response we set
	err = client.TriggerModeScriptEventArray("Maniaplanet.Pause.SetActive", []string{"true"})
	if err != nil {
		log.Fatal("Failed to trigger script event:", err)
	}
	fmt.Println("Successfully triggered script event!")

	// Test another method call
	err = client.RestartMap()
	if err != nil {
		log.Fatal("Failed to restart map:", err)
	}
	fmt.Println("Successfully called RestartMap!")

	fmt.Println("Example completed successfully!")
}
