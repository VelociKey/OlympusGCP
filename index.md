# OlympusGCP-FinOps Index

This workspace houses the cost estimation and financial governance tools for the Olympus fleet.

## Actors (10000)
- [FinOpsManager](./10000-Autonomous-Actors/900-FinOpsManager/): Go-based ConnectRPC backend for real-time cost heuristics.

## Bridges (20000)
- [FinOpsBridge](./30000-Context-Bridges/900-FinOpsBridge/): MCP-Native frontend exposing cost prediction tools (`finops_estimate_cost`).

## Contracts (40000)
- [Proto Contracts](./50000-Communication-Contracts/proto/): Protobuf definitions for financial services.
