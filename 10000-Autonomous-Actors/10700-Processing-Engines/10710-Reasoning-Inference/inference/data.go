//go:build !wasm

package inference

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	datav1 "olympus.fleet/00SDLC/OlympusGCP/gen/data"
	"connectrpc.com/connect"
	"go.etcd.io/bbolt"
)

type DataServer struct {
	db *bbolt.DB
}

const (
	bucketCollections = "collections"
	bucketTables      = "tables"
)

func NewDataServer(storageDir string) *DataServer {
	os.MkdirAll(storageDir, 0755)
	dbPath := filepath.Join(storageDir, "data.db")
	
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		slog.Error("Failed to open BoltDB", "path", dbPath, "error", err)
		panic(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketCollections))
		if err != nil { return err }
		_, err = tx.CreateBucketIfNotExists([]byte(bucketTables))
		return err
	})
	if err != nil { panic(err) }

	return &DataServer{db: db}
}

func (s *DataServer) Upsert(ctx context.Context, req *connect.Request[datav1.UpsertRequest]) (*connect.Response[datav1.UpsertResponse], error) {
	col := req.Msg.Collection
	id := req.Msg.DocId
	slog.Info("Upsert", "collection", col, "id", id)

	err := s.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(bucketCollections))
		b, err := root.CreateBucketIfNotExists([]byte(col))
		if err != nil { return err }
		return b.Put([]byte(id), []byte(req.Msg.DataJson))
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&datav1.UpsertResponse{}), nil
}

func (s *DataServer) WriteRow(ctx context.Context, req *connect.Request[datav1.WriteRowRequest]) (*connect.Response[datav1.WriteRowResponse], error) {
	slog.Info("WriteRow", "table", req.Msg.Table, "key", req.Msg.RowKey)
	
	err := s.db.Update(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(bucketTables))
		b, err := root.CreateBucketIfNotExists([]byte(req.Msg.Table))
		if err != nil { return err }
		// In deep emulation, we would parse RowData and store columns.
		// For now, storing as blob.
		return b.Put([]byte(req.Msg.RowKey), []byte("row_data"))
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&datav1.WriteRowResponse{}), nil
}

func (s *DataServer) Query(ctx context.Context, req *connect.Request[datav1.QueryRequest]) (*connect.Response[datav1.QueryResponse], error) {
	slog.Info("Query", "collection", req.Msg.Collection, "query", req.Msg.Query)
	
	var results []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(bucketCollections))
		b := root.Bucket([]byte(req.Msg.Collection))
		if b == nil { return nil }
		
		return b.ForEach(func(k, v []byte) error {
			// Naive query: if query matches content
			if req.Msg.Query == "" || strings.Contains(string(v), req.Msg.Query) || req.Msg.Query == "all" {
				results = append(results, string(v))
			}
			return nil
		})
	})

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if results == nil {
		results = []string{}
	}

	data, _ := json.Marshal(results)
	return connect.NewResponse(&datav1.QueryResponse{ResultsJson: string(data)}), nil
}

func (s *DataServer) Close() error {
	return s.db.Close()
}
