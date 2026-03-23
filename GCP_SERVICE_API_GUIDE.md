# Google Cloud Platform: Comprehensive Service & API Reference

This manifest provides a high-level mapping of GCP services and their associated API identifiers. Use this guide to identify candidates for local-dev emulators and ConnectRPC stubs.

## 1. Compute & Serverless
| Service | API Identifier | Latest Stable | Description |
| :--- | :--- | :---: | :--- |
| **Compute Engine** | `compute.googleapis.com` | `v1` | Virtual Machines and Managed Instance Groups. |
| **Cloud Run** | `run.googleapis.com` | `v2` | Managed Knative-based container execution. |
| **App Engine** | `appengine.googleapis.com` | `v1` | PaaS for web applications. |
| **GKE Hub** | `gkehub.googleapis.com` | `v1` | Fleet management for Kubernetes clusters. |

## 2. Information Storage & Databases
| Service | API Identifier | Latest Stable | Description |
| :--- | :--- | :---: | :--- |
| **Cloud Storage** | `storage.googleapis.com` | `v1` | Unstructured Object (Blob) storage. |
| **Cloud Firestore** | `firestore.googleapis.com` | `v1` | Document-oriented NoSQL database. |
| **Cloud Spanner** | `spanner.googleapis.com` | `v1` | Horizontally scalable relational database. |
| **Cloud SQL** | `sqladmin.googleapis.com` | `v1` | Managed MySQL, PostGreSQL, and SQL Server. |
| **Bigtable** | `bigtableadmin.googleapis.com` | `v2` | High-performance wide-column NoSQL. |

## 3. Big Data & Messaging
| Service | API Identifier | Latest Stable | Description |
| :--- | :--- | :---: | :--- |
| **BigQuery** | `bigquery.googleapis.com` | `v2` | Data warehouse and analytics. |
| **Cloud Pub/Sub** | `pubsub.googleapis.com` | `v1` | Asynchronous messaging (Topics/Subscriptions). |
| **Cloud Dataflow** | `dataflow.googleapis.com` | `v1` | Stream and batch data processing. |
| **Cloud Dataproc** | `dataproc.googleapis.com` | `v1` | Managed Apache Spark and Hadoop. |

## 4. Artificial Intelligence & Machine Learning
| Service | API Identifier | Latest Stable | Description |
| :--- | :--- | :---: | :--- |
| **Vertex AI** | `aiplatform.googleapis.com` | `v1` | Unified AI platform for training and inference. |
| **Natural Language**| `language.googleapis.com` | `v2` | Text analysis and sentiment extraction. |
| **Vision AI** | `vision.googleapis.com` | `v1` | Image recognition and OCR. |
| **Speech-to-Text** | `speech.googleapis.com` | `v2` | Audio transcription and recognition. |

## 5. Security & Management
| Service | API Identifier | Latest Stable | Description |
| :--- | :--- | :---: | :--- |
| **IAM** | `iam.googleapis.com` | `v1` | Identity and Access Management. |
| **Secret Manager** | `secretmanager.googleapis.com` | `v1` | Management of confidential data/keys. |
| **Cloud Logging** | `logging.googleapis.com` | `v2` | Centralized log ingestion and analytics. |
| **Cloud Monitoring** | `monitoring.googleapis.com` | `v3` | System metrics and observability. |
| **Resource Manager**| `cloudresourcemanager.googleapis.com`| `v3` | Project/Folder hierarchy management. |

## 6. Events & Orchestration
| Service | API Identifier | Latest Stable | Description |
| :--- | :--- | :---: | :--- |
| **Eventarc** | `eventarc.googleapis.com` | `v1` | Event-driven architecture and triggers. |
| **Cloud Tasks** | `cloudtasks.googleapis.com` | `v2` | Distributed task queues and retries. |
| **Cloud Workflows** | `workflows.googleapis.com` | `v1` | Serverless orchestration of services. |

---
*Reference generated from Google APIs Discovery Service | 2026-03-22*
