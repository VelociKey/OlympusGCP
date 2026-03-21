package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/mark3labs/mcp-go/mcp"

	"olympus.fleet/00SDLC/OlympusGCP-Data/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/data/datav1connect"
	datav1 "olympus.fleet/00SDLC/OlympusGCP-Data/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/data"
	"olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/140-MCPBridge"
)

func main() {
	s := mcpbridge.NewBridgeServer("OlympusDataBridge", "1.0.0")

	client := datav1connect.NewDataServiceClient(
		http.DefaultClient,
		"http://localhost:8093",
	)

	s.AddTool(mcp.NewTool("data_upsert",
		mcp.WithDescription("Upsert a document into a collection. Args: {collection: string, doc_id: string, data_json: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		collection, _ := m["collection"].(string)
		docID, _ := m["doc_id"].(string)
		dataJSON, _ := m["data_json"].(string)

		_, err = client.Upsert(ctx, connect.NewRequest(&datav1.UpsertRequest{
			Collection: collection,
			DocId:      docID,
			DataJson:   dataJSON,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Document '%s' upserted into collection '%s'.", docID, collection)), nil
	})

	s.AddTool(mcp.NewTool("data_query",
		mcp.WithDescription("Query documents in a collection. Args: {collection: string, query: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		collection, _ := m["collection"].(string)
		query, _ := m["query"].(string)

		resp, err := client.Query(ctx, connect.NewRequest(&datav1.QueryRequest{
			Collection: collection,
			Query:      query,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Query results for '%s': %s", collection, resp.Msg.ResultsJson)), nil
	})

	s.Run()
}
