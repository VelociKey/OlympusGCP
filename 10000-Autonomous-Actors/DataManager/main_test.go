package main

import (
	"context"
	"testing"

	datav1 "olympus.fleet/00SDLC/OlympusGCP-Data/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/data"
	"olympus.fleet/00SDLC/OlympusGCP-Data/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestDataServer(t *testing.T) {
	tempDir := t.TempDir()
	server := inference.NewDataServer(tempDir)
	defer server.Close()
	ctx := context.Background()

	// Test Upsert
	upsertReq := connect.NewRequest(&datav1.UpsertRequest{
		Collection: "users",
		DocId:      "jdoe",
		DataJson:   `{"name": "John Doe"}`,
	})
	_, err := server.Upsert(ctx, upsertReq)
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// Test WriteRow
	writeReq := connect.NewRequest(&datav1.WriteRowRequest{
		Table:  "logs",
		RowKey: "2026-02-23",
	})
	_, err = server.WriteRow(ctx, writeReq)
	if err != nil {
		t.Fatalf("WriteRow failed: %v", err)
	}

	// Test Query
	queryReq := connect.NewRequest(&datav1.QueryRequest{
		Collection: "users",
		Query:      "name == 'John Doe'",
	})
	queryRes, err := server.Query(ctx, queryReq)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if queryRes.Msg.ResultsJson != "[]" {
		t.Errorf("Expected empty results, got %s", queryRes.Msg.ResultsJson)
	}
}
