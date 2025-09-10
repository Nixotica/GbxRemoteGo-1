# GbxRemoteGo

A package for interacting with the XML-RPC protocol of Trackmania servers.

## Installation

Install GbxRemoteGo with Go's package manager

```bash
go get github.com/MRegterschot/GbxRemoteGo
```

## Usage/Examples

```go
package main

import (
	"fmt"
	"os"

	"github.com/MRegterschot/GbxRemoteGo/events"
	. "github.com/MRegterschot/GbxRemoteGo/gbxclient"
)

func main() {
	// Create a new GbxClient
	client := NewGbxClient(Options{})

	// Register event handlers
	onConnectionChan := make(chan any)
	client.Events.On("connect", onConnectionChan)
	go handleConnect(onConnectionChan)

	onDisconnectChan := make(chan any)
	client.Events.On("disconnect", onDisconnectChan)
	go handleDisconnect(onDisconnectChan)

	// Connect to the server
	if err := client.Connect("127.0.0.1", 5000); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := client.SetApiVersion("2023-04-24"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := client.EnableCallbacks(true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := client.Authenticate("SuperAdmin", "SuperAdmin"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Register gbx callback handlers
	client.OnPlayerConnect = append(client.OnPlayerConnect, func(client *GbxClient, args events.PlayerConnectEventArgs) {
		fmt.Println("Player connected:", args.Login)
	})

	client.OnPlayerCheckpoint = append(client.OnPlayerCheckpoint, func(client *GbxClient, args events.PlayerCheckpointEventArgs) {
		fmt.Println("Player checkpoint:", args)
	})

	client.OnAnyCallback = append(client.OnAnyCallback, func(client *GbxClient, args CallbackEventArgs) {
		fmt.Println("Any callback:", args)
	})

	select {}
}

func handleConnect(eventChan chan any) {
	for {
		select {
		case event := <-eventChan:
			if connected, ok := event.(bool); ok {
				if connected {
					fmt.Println("Connected")
				} else {
					fmt.Println("Not Connected")
				}
			} else {
				fmt.Println("Invalid event type for connect.")
			}
		}
	}
}

func handleDisconnect(eventChan chan any) {
	for {
		select {
		case event := <-eventChan:
			if msg, ok := event.(string); ok {
				fmt.Println(msg)
			} else {
				fmt.Println("Invalid event type for disconnect.")
			}
		}
	}
}
```

## Mock Server

GbxRemoteGo includes a mock server for testing and development purposes. The mock server implements the same XML-RPC protocol as real Trackmania servers, allowing you to test your applications without needing a running game server.

### Basic Usage

```go
package main

import (
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
    
    if err := mock.Start(); err != nil {
        log.Fatal("Failed to start mock server:", err)
    }
    defer mock.Stop()
    
    // Set custom responses for specific methods
    mock.SetResponse("TriggerModeScriptEventArray", "Pause activated")
    mock.SetResponse("GetVersion", map[string]interface{}{
        "Name":    "MyMockServer",
        "Version": "1.0.0",
    })
    
    // Connect client to mock server
    client := gbxclient.NewGbxClient("127.0.0.1", 5000, gbxclient.Options{})
    client.Connect()
    client.SetApiVersion("2023-04-24")
    client.Authenticate("SuperAdmin", "SuperAdmin")
    
    // Use client normally - calls go to mock server
    version, _ := client.GetVersion()
    println("Connected to:", version.Name)
}
```

### Features

- **Full XML-RPC compatibility** with existing client code
- **Custom response injection** for specific method calls
- **Default responses** for common methods (authentication, version info, etc.)
- **Easy testing setup** for unit tests and integration tests
- **No external dependencies** beyond the existing codebase

### Default Supported Methods

The mock server provides default responses for:
- `SetApiVersion`, `EnableCallbacks`, `Authenticate` - Always return success
- `GetVersion`, `GetStatus`, `GetSystemInfo` - Return mock server information
- `TriggerModeScriptEvent`, `TriggerModeScriptEventArray` - Accept calls without error
- `RestartMap`, `NextMap` - Accept calls without error

For any other method calls, the mock server returns a generic success response unless you've set a custom response using `SetResponse()`.

