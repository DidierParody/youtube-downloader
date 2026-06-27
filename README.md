# Youtube Downloader

Plataforma de gestion y descarga inteligente de videos de YouTube construida con arquitectura distribuida.

## Arquitectura

- **Frontend**: SvelteKit + Cloudflare Pages
- **Backend**: Go (Fiber)
- **Base de Datos OLTP**: PostgreSQL (Neon)
- **Vector Search**: pgvector
- **Cache**: Redis
- **Bus de Eventos**: Redpanda (Compatible Kafka)
- **Workers**: Go
- **Storage**: Cloudflare R2 (S3 Compatible)
- **Lakehouse**: Apache Iceberg + Parquet
- **Motor Analitico**: DuckDB

## Modelo de Datos

El modelo de datos esta completamente documentado en la carpeta `database/`.

### Entregables

1. **MODEL.md**: Modelo OLTP completamente documentado con explicacion de cada tabla, proposito, relaciones y decisiones de diseño.
2. **MER.mmd**: Diagrama ER en Mermaid compatible con Excalidraw (entidades, atributos, relaciones, cardinalidades).
3. **schema.sql**: DDL PostgreSQL completo y ejecutable en Neon, incluyendo extensiones, tablas, indices, triggers, particionamiento y comentarios.
4. **OPTIMIZATION.md**: Documento de optimizacion con explicacion de indices, consultas optimizadas, recomendaciones de VACUUM/ANALYZE, pgvector y mejoras futuras.

### Convenciones del Modelo

- **Claves primarias**: UUID v7 (generados en Go)
- **Timestamps**: Siempre TIMESTAMPTZ para consistencia en sistemas distribuidos
- **Patron de auditoria**: `created_at` + `updated_at` + trigger automatico para `updated_at`
- **Estados y tipos**: ENUMs para conjuntos estables de valores
- **Particionamiento**: `AuditEvent` particionado mensualmente por rango en `occurred_at`

### Entidades Principales

- **Usuario**: Cuentas registradas con cuotas de almacenamiento
- **Video**: Contenido logico de YouTube
- **Archivo**: Representacion fisica con deduplicacion (sha256 + reference_count)
- **Descarga**: Solicitudes de descarga con estado del pipeline
- **Embedding**: Vectores para busqueda semantica con pgvector
- **WorkerExecution**: Registro de ejecuciones del pipeline
- **AuditEvent**: Auditoria completa de todos los eventos del sistema

## Instalacion

Instrucciones de instalacion y desarrollo se anadiran proximamente.

## Licencia

Este proyecto es open source y esta disponible para todos.
