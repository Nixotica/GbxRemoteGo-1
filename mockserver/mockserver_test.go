package mockserver

import (
	"testing"
	"time"

	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
)

func TestMockServer(t *testing.T) {
	// Create and start mock server
	mock := New(Config{
		Host: "127.0.0.1",
		Port: 5001, // Use different port for tests
	})

	// Set a custom response
	mock.SetResponse("GetVersion", map[string]interface{}{
		"Name":    "TestServer",
		"Version": "1.0.0-test",
	})

	if err := mock.Start(); err != nil {
		t.Fatal("Failed to start mock server:", err)
	}
	defer mock.Stop()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Create client and connect
	client := gbxclient.NewGbxClient("127.0.0.1", 5001, gbxclient.Options{})
	if err := client.Connect(); err != nil {
		t.Fatal("Failed to connect to mock server:", err)
	}
	defer client.Disconnect()

	// Test authentication flow
	if err := client.SetApiVersion("2023-04-24"); err != nil {
		t.Fatal("Failed to set API version:", err)
	}

	if err := client.EnableCallbacks(true); err != nil {
		t.Fatal("Failed to enable callbacks:", err)
	}

	if err := client.Authenticate("test", "test"); err != nil {
		t.Fatal("Failed to authenticate:", err)
	}

	// Test custom response
	version, err := client.GetVersion()
	if err != nil {
		t.Fatal("Failed to get version:", err)
	}

	if version.Name != "TestServer" {
		t.Errorf("Expected Name='TestServer', got '%s'", version.Name)
	}

	if version.Version != "1.0.0-test" {
		t.Errorf("Expected Version='1.0.0-test', got '%s'", version.Version)
	}

	// Test default response
	status, err := client.GetStatus()
	if err != nil {
		t.Fatal("Failed to get status:", err)
	}

	if status.Code != 4 {
		t.Errorf("Expected status code 4, got %d", status.Code)
	}

	// Test script event methods
	err = client.TriggerModeScriptEvent("Trackmania.Pause.SetActive", "true")
	if err != nil {
		t.Fatal("Failed to trigger script event:", err)
	}

	err = client.TriggerModeScriptEventArray("Trackmania.ForceEndRound", []string{})
	if err != nil {
		t.Fatal("Failed to trigger script event array:", err)
	}
}

func TestMockServerCustomResponses(t *testing.T) {
	mock := New(Config{
		Host: "127.0.0.1",
		Port: 5002,
	})

	// Test different response types
	mock.SetResponse("TestString", "Hello World")
	mock.SetResponse("TestInt", 42)
	mock.SetResponse("TestBool", true)
	mock.SetResponse("TestStruct", map[string]interface{}{
		"Field1": "value1",
		"Field2": 123,
		"Field3": false,
	})

	if err := mock.Start(); err != nil {
		t.Fatal("Failed to start mock server:", err)
	}
	defer mock.Stop()

	time.Sleep(100 * time.Millisecond)

	client := gbxclient.NewGbxClient("127.0.0.1", 5002, gbxclient.Options{})
	if err := client.Connect(); err != nil {
		t.Fatal("Failed to connect:", err)
	}
	defer client.Disconnect()

	// Basic auth
	client.SetApiVersion("2023-04-24")
	client.EnableCallbacks(true)
	client.Authenticate("test", "test")

	// These would be tested if we had generic Call method exposed
	// For now, we test that the mock server handles unknown methods gracefully
	err := client.TriggerModeScriptEvent("TestString", "")
	if err != nil {
		t.Fatal("Mock server should handle unknown methods gracefully:", err)
	}
}
