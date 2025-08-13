package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/v-saba/bazel-tutorial/proto/gen/go/telemetry_server/v1"

	"github.com/v-saba/bazel-tutorial/common"
)

// server implements the gRPC TelemetryServiceServer interface
type server struct {
	pb.UnimplementedTelemetryServiceServer
}

// QueryTelemetry implements the gRPC TelemetryServiceServer interface
func (s *server) QueryTelemetry(ctx context.Context, req *pb.TelemetryRequest) (*pb.TelemetryResponse, error) {
	log.Printf("Received gRPC telemetry request: %v", req)

	// Create response based on the request type
	response := &pb.TelemetryResponse{
		TelemetryType: req.TelemetryType,
		Timestamp:     time.Now().Unix(),
	}

	switch req.TelemetryType {
	case pb.TelemetryType_TELEMETRY_TYPE_HEARTBEAT:
		response.TelemetryData = &pb.TelemetryResponse_HeartbeatTelemetry{
			HeartbeatTelemetry: &pb.HeartbeatTelemetry{},
		}
	case pb.TelemetryType_TELEMETRY_TYPE_LOG:
		response.TelemetryData = &pb.TelemetryResponse_LogTelemetry{
			LogTelemetry: &pb.LogTelemetry{
				LogData: "Sample log data from gRPC server",
			},
		}
	case pb.TelemetryType_TELEMETRY_TYPE_CPU_USAGE:
		response.TelemetryData = &pb.TelemetryResponse_CpuUsageTelemetry{
			CpuUsageTelemetry: &pb.CpuUsageTelemetry{
				CpuUsage: 42.5, // Mock CPU usage
			},
		}
	default:
		log.Printf("Unknown telemetry type: %v", req.TelemetryType)
		return nil, fmt.Errorf("unknown telemetry type: %v", req.TelemetryType)
	}

	log.Printf("Sending gRPC telemetry response: %v", response)
	return response, nil
}

// ProcessTelemetryRequest processes a telemetry request and returns a response (legacy method)
func (s *server) ProcessTelemetryRequest(req *pb.TelemetryRequest) *pb.TelemetryResponse {
	log.Printf("Processing telemetry request: %v", req)

	// Create response based on the request type
	response := &pb.TelemetryResponse{
		TelemetryType: req.TelemetryType,
		Timestamp:     time.Now().Unix(),
	}

	switch req.TelemetryType {
	case pb.TelemetryType_TELEMETRY_TYPE_HEARTBEAT:
		response.TelemetryData = &pb.TelemetryResponse_HeartbeatTelemetry{
			HeartbeatTelemetry: &pb.HeartbeatTelemetry{},
		}
	case pb.TelemetryType_TELEMETRY_TYPE_LOG:
		response.TelemetryData = &pb.TelemetryResponse_LogTelemetry{
			LogTelemetry: &pb.LogTelemetry{
				LogData: "Sample log data from server",
			},
		}
	case pb.TelemetryType_TELEMETRY_TYPE_CPU_USAGE:
		response.TelemetryData = &pb.TelemetryResponse_CpuUsageTelemetry{
			CpuUsageTelemetry: &pb.CpuUsageTelemetry{
				CpuUsage: 42.5, // Mock CPU usage
			},
		}
	default:
		log.Printf("Unknown telemetry type: %v", req.TelemetryType)
		response = nil
	}

	if response != nil {
		log.Printf("Generated telemetry response: %v", response)
	}
	return response
}

func testMessages() {
	// Test creating basic message types
	req := &pb.TelemetryRequest{
		TelemetryType: pb.TelemetryType_TELEMETRY_TYPE_HEARTBEAT,
	}

	resp := &pb.TelemetryResponse{
		TelemetryType: pb.TelemetryType_TELEMETRY_TYPE_HEARTBEAT,
		TelemetryData: &pb.TelemetryResponse_HeartbeatTelemetry{
			HeartbeatTelemetry: &pb.HeartbeatTelemetry{},
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("Test Request: %v", req)
	log.Printf("Test Response: %v", resp)
}

func main() {
	newUUID := common.GenerateUUIDStr()
	log.Printf("New UUID: %v", newUUID)

	log.Printf("Telemetry server starting...")

	// Test basic proto messages first
	testMessages()

	// Create a server instance
	srv := &server{}

	// Test different telemetry types with both methods
	telemetryTypes := []pb.TelemetryType{
		pb.TelemetryType_TELEMETRY_TYPE_HEARTBEAT,
		pb.TelemetryType_TELEMETRY_TYPE_LOG,
		pb.TelemetryType_TELEMETRY_TYPE_CPU_USAGE,
	}

	log.Printf("Testing legacy method...")
	for _, telType := range telemetryTypes {
		req := &pb.TelemetryRequest{
			TelemetryType: telType,
		}

		resp := srv.ProcessTelemetryRequest(req)
		if resp != nil {
			log.Printf("Legacy: Successfully processed %v request", telType)
		} else {
			log.Printf("Legacy: Failed to process %v request", telType)
		}
	}

	log.Printf("Testing gRPC method...")
	for _, telType := range telemetryTypes {
		req := &pb.TelemetryRequest{
			TelemetryType: telType,
		}

		resp, err := srv.QueryTelemetry(context.Background(), req)
		if err != nil {
			log.Printf("gRPC: Failed to process %v request: %v", telType, err)
		} else if resp != nil {
			log.Printf("gRPC: Successfully processed %v request", telType)
		}
	}

	log.Printf("Server test completed successfully!")
}
