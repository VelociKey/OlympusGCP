package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/mark3labs/mcp-go/mcp"

	"olympus.fleet/00SDLC/OlympusGCP-Vault/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/vault/vaultv1connect"
	vaultv1 "olympus.fleet/00SDLC/OlympusGCP-Vault/40000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/gen/v1/vault"
	"olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/140-MCPBridge"
)

func main() {
	s := mcpbridge.NewBridgeServer("OlympusVaultBridge", "1.0.0")

	client := vaultv1connect.NewVaultServiceClient(
		http.DefaultClient,
		"http://localhost:8092",
	)

	s.AddTool(mcp.NewTool("vault_write",
		mcp.WithDescription("Write a secret to the local Vault. Args: {key: string, value: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		key, _ := m["key"].(string)
		value, _ := m["value"].(string)

		_, err = client.VaultWrite(ctx, connect.NewRequest(&vaultv1.VaultWriteRequest{
			Key:   key,
			Value: value,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Secret '%s' successfully written to Vault.", key)), nil
	})

	s.AddTool(mcp.NewTool("vault_read",
		mcp.WithDescription("Read a secret from the local Vault. Args: {key: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		key, _ := m["key"].(string)

		resp, err := client.VaultRead(ctx, connect.NewRequest(&vaultv1.VaultReadRequest{
			Key: key,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Secret '%s': %s", key, resp.Msg.Value)), nil
	})

	s.Run()
}
