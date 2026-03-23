//go:build !wasm

package inference

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	vaultv1 "olympus.fleet/00SDLC/OlympusGCP/gen/vault"
	"olympus.fleet/00SDLC/Olympus2/90000-Enablement-Labs/90200-Logic-Libraries/170-Policy"
	"connectrpc.com/connect"
	"go.etcd.io/bbolt"
)

type VaultServer struct {
	mu     sync.RWMutex
	db     *bbolt.DB
	pe     *policy.Evaluator
}

const (
	bucketSecrets = "secrets"
	bucketHistory = "history"
)

func NewVaultServer(storageDir string) *VaultServer {
	os.MkdirAll(storageDir, 0755)
	dbPath := filepath.Join(storageDir, "vault.db")
	
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		slog.Error("Failed to open BoltDB", "path", dbPath, "error", err)
		panic(err)
	}

	// Initialize buckets
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketSecrets))
		if err != nil { return err }
		_, err = tx.CreateBucketIfNotExists([]byte(bucketHistory))
		return err
	})
	if err != nil {
		slog.Error("Failed to init buckets", "error", err)
		panic(err)
	}

	s := &VaultServer{
		db: db,
		pe: &policy.Evaluator{},
	}
	
	// Load PBAC policy
	cwd, _ := os.Getwd()
	slog.Info("VaultServer initializing", "cwd", cwd)
	policyPaths := []string{
		"../../olympus.fleet/20POSI/George/C0100-Configuration-Registry/POLICY.jebnf",
		"olympus.fleet/20POSI/George/C0100-Configuration-Registry/POLICY.jebnf",
		"../../../../../olympus.fleet/20POSI/George/C0100-Configuration-Registry/POLICY.jebnf",
	}
	for _, p := range policyPaths {
		if err := s.pe.Load(p); err == nil {
			slog.Info("Successfully loaded PBAC policy", "path", p)
			break
		}
	}

	return s
}

func (s *VaultServer) VaultWrite(ctx context.Context, req *connect.Request[vaultv1.VaultWriteRequest]) (*connect.Response[vaultv1.VaultWriteResponse], error) {
	key := req.Msg.Key
	val := req.Msg.Value
	slog.Info("VaultWrite", "key", key)

	var version int32
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketSecrets))
		if err := b.Put([]byte(key), []byte(val)); err != nil { return err }

		h := tx.Bucket([]byte(bucketHistory))
		histData := h.Get([]byte(key))
		var versions []string
		if histData != nil {
			json.Unmarshal(histData, &versions)
		}
		versions = append(versions, val)
		version = int32(len(versions))
		
		newData, _ := json.Marshal(versions)
		return h.Put([]byte(key), newData)
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&vaultv1.VaultWriteResponse{Version: version}), nil
}

func (s *VaultServer) VaultRead(ctx context.Context, req *connect.Request[vaultv1.VaultReadRequest]) (*connect.Response[vaultv1.VaultReadResponse], error) {
	slog.Info("VaultRead", "key", req.Msg.Key)
	
	var val string
	var version int32
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketSecrets))
		v := b.Get([]byte(req.Msg.Key))
		if v == nil {
			return fmt.Errorf("secret not found: %s", req.Msg.Key)
		}
		val = string(v)

		h := tx.Bucket([]byte(bucketHistory))
		histData := h.Get([]byte(req.Msg.Key))
		var versions []string
		json.Unmarshal(histData, &versions)
		version = int32(len(versions))
		return nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&vaultv1.VaultReadResponse{Value: val, Version: version}), nil
}

func (s *VaultServer) GetSecretVersion(ctx context.Context, req *connect.Request[vaultv1.GetSecretVersionRequest]) (*connect.Response[vaultv1.VaultReadResponse], error) {
	slog.Info("GetSecretVersion", "key", req.Msg.Key, "version", req.Msg.Version)
	
	var val string
	err := s.db.View(func(tx *bbolt.Tx) error {
		h := tx.Bucket([]byte(bucketHistory))
		histData := h.Get([]byte(req.Msg.Key))
		if histData == nil {
			return fmt.Errorf("secret not found: %s", req.Msg.Key)
		}
		
		var versions []string
		json.Unmarshal(histData, &versions)
		
		idx := int(req.Msg.Version) - 1
		if idx < 0 || idx >= len(versions) {
			return fmt.Errorf("version %d not found for %s", req.Msg.Version, req.Msg.Key)
		}
		val = versions[idx]
		return nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&vaultv1.VaultReadResponse{Value: val, Version: req.Msg.Version}), nil
}

func (s *VaultServer) ListSecretVersions(ctx context.Context, req *connect.Request[vaultv1.ListSecretVersionsRequest]) (*connect.Response[vaultv1.ListSecretVersionsResponse], error) {
	slog.Info("ListSecretVersions", "key", req.Msg.Key)
	
	var resVersions []int32
	err := s.db.View(func(tx *bbolt.Tx) error {
		h := tx.Bucket([]byte(bucketHistory))
		histData := h.Get([]byte(req.Msg.Key))
		if histData == nil {
			return fmt.Errorf("secret not found: %s", req.Msg.Key)
		}
		
		var versions []string
		json.Unmarshal(histData, &versions)
		
		for i := range versions {
			resVersions = append(resVersions, int32(i+1))
		}
		return nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&vaultv1.ListSecretVersionsResponse{Versions: resVersions}), nil
}

func (s *VaultServer) ListSecrets(ctx context.Context, req *connect.Request[vaultv1.ListSecretsRequest]) (*connect.Response[vaultv1.ListSecretsResponse], error) {
	slog.Info("ListSecrets", "prefix", req.Msg.Prefix)
	
	var keys []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketSecrets))
		c := b.Cursor()
		
		prefix := []byte(req.Msg.Prefix)
		for k, _ := c.Seek(prefix); k != nil && strings.HasPrefix(string(k), string(prefix)); k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&vaultv1.ListSecretsResponse{Keys: keys}), nil
}

func (s *VaultServer) TestIAMPolicy(ctx context.Context, req *connect.Request[vaultv1.TestIAMPolicyRequest]) (*connect.Response[vaultv1.TestIAMPolicyResponse], error) {
	slog.Info("TestIAMPolicy", "identity", req.Msg.Identity, "action", req.Msg.Action, "resource", req.Msg.Resource)
	
	allowed, reason := s.pe.Authorize("George", req.Msg.Identity, "*")
	
	return connect.NewResponse(&vaultv1.TestIAMPolicyResponse{
		Allowed: allowed,
		Reason:  reason,
	}), nil
}

func (s *VaultServer) Close() error {
	return s.db.Close()
}
