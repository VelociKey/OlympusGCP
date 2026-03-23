package inference

import (
	"context"
	"strings"
	"testing"

	datav1 "olympus.fleet/00SDLC/OlympusGCP/gen/data"
	"connectrpc.com/connect"
)

func TestDataServer_CoverageExpansion(t *testing.T) {
	tempDir := t.TempDir()
	server := NewDataServer(tempDir)
	defer server.Close()
	ctx := context.Background()

	// 1. Test Upsert
	docJSON := `{"val":1}`
	_, err := server.Upsert(ctx, connect.NewRequest(&datav1.UpsertRequest{
		Collection: "test",
		DocId:      "doc1",
		DataJson:   docJSON,
	}))
	if err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	// 2. Test Query
	res, err := server.Query(ctx, connect.NewRequest(&datav1.QueryRequest{
		Collection: "test",
		Query:      "all",
	}))
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if !strings.Contains(res.Msg.ResultsJson, `{\"val\":1}`) {
		t.Errorf("Expected results to contain escaped json, got %s", res.Msg.ResultsJson)
	}

	// 3. Test WriteRow
	_, err = server.WriteRow(ctx, connect.NewRequest(&datav1.WriteRowRequest{
		Table:  "table1",
		RowKey: "row1",
	}))
	if err != nil {
		t.Error("WriteRow failed")
	}
}
