package inference

import (
	"context"
	"os"
	"testing"

	vaultv1 "olympus.fleet/00SDLC/OlympusGCP/gen/vault"
	"connectrpc.com/connect"
)

func TestVaultServer_CoverageExpansion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "vault-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	server := NewVaultServer(tempDir)
	ctx := context.Background()

	// 1. Test Read Non-Existent Key
	_, err = server.VaultRead(ctx, connect.NewRequest(&vaultv1.VaultReadRequest{Key: "missing"}))
	if err == nil {
		t.Error("Expected error for missing key, got nil")
	}

	// 2. Test List Versions Non-Existent Key
	_, err = server.ListSecretVersions(ctx, connect.NewRequest(&vaultv1.ListSecretVersionsRequest{Key: "missing"}))
	if err == nil {
		t.Error("Expected error for missing key versions, got nil")
	}

	// 3. Test Get Secret Version Out of Range
	server.VaultWrite(ctx, connect.NewRequest(&vaultv1.VaultWriteRequest{Key: "k1", Value: "olympus.fleet/00SDLC/Olympus2/50000-Communication-Contracts/40400-Protocol-Synthetics/connect-rpc/olympus/v1"}))
	_, err = server.GetSecretVersion(ctx, connect.NewRequest(&vaultv1.GetSecretVersionRequest{Key: "k1", Version: 100}))
	if err == nil {
		t.Error("Expected error for out-of-range version, got nil")
	}

	// 4. Test List Secrets with Prefix
	server.VaultWrite(ctx, connect.NewRequest(&vaultv1.VaultWriteRequest{Key: "app/db", Value: "pass"}))
	server.VaultWrite(ctx, connect.NewRequest(&vaultv1.VaultWriteRequest{Key: "app/api", Value: "key"}))
	server.VaultWrite(ctx, connect.NewRequest(&vaultv1.VaultWriteRequest{Key: "other/stuff", Value: "val"}))

	res, err := server.ListSecrets(ctx, connect.NewRequest(&vaultv1.ListSecretsRequest{Prefix: "app/"}))
	if err != nil {
		t.Fatalf("ListSecrets failed: %v", err)
	}
	if len(res.Msg.Keys) != 2 {
		t.Errorf("Expected 2 keys with prefix 'app/', got %d", len(res.Msg.Keys))
	}

	// 5. Test IAM Policy (Deep Emulation)
	policyRes, _ := server.TestIAMPolicy(ctx, connect.NewRequest(&vaultv1.TestIAMPolicyRequest{Identity: "user-123", Action: "read"}))
	if !policyRes.Msg.Allowed {
		t.Error("Expected IAM policy to be allowed for 'user-123'")
	}

	failRes, _ := server.TestIAMPolicy(ctx, connect.NewRequest(&vaultv1.TestIAMPolicyRequest{Identity: "unauthorized", Action: "read"}))
	if failRes.Msg.Allowed {
		t.Error("Expected IAM policy to be denied for 'unauthorized'")
	}
}
