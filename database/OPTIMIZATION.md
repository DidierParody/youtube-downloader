# Documento de Optimizacion - Youtube Downloader

## Introduccion

Este documento explica las decisiones de indexacion, optimizacion y mantenimiento para el modelo de datos OLTP de la plataforma de descarga de videos.

## Indices por Tabla

### Usuario

- **idx_usuario_email (email)**: Optimiza autenticacion y busquedas por email. Unico por constraint UNIQUE, indice automatico incluido.
- **idx_usuario_username (username)**: Optimiza busquedas por nombre de usuario. Unico por constraint UNIQUE, indice automatico incluido.

**Consultas optimizadas**:
- Login de usuario: `SELECT * FROM "Usuario" WHERE email = $1`
- Perfil publico: `SELECT * FROM "Usuario" WHERE username = $1`

### Video

- **idx_video_youtube_id (youtube_video_id)**: Unico por constraint UNIQUE, evita duplicados y acelera busquedas por ID de YouTube.
- **idx_video_title (title)**: Facilita busquedas de texto parcial (LIKE) en titulos de videos.

**Consultas optimizadas**:
- Verificar existencia de video: `SELECT id FROM "Video" WHERE youtube_video_id = $1`
- Busqueda por titulo: `SELECT * FROM "Video" WHERE title ILIKE '%term%'`

### Archivo

- **idx_archivo_video_id (video_id)**: Acelera obtencion de archivos asociados a un video.
- **idx_archivo_sha256 (sha256)**: Unico por constraint UNIQUE, esencial para deduplicacion.
- **idx_archivo_object_key (object_key)**: Necesario para operaciones de storage (R2/S3).

**Consultas optimizadas**:
- Obtener archivos de un video: `SELECT * FROM "Archivo" WHERE video_id = $1`
- Deduplicacion: `SELECT id FROM "Archivo" WHERE sha256 = $1`
- Acceso a storage: `SELECT object_key FROM "Archivo" WHERE id = $1`

### Descarga

- **idx_descarga_user_id (user_id)**: Historial de descargas de un usuario.
- **idx_descarga_video_id (video_id)**: Quien descargo un video especifico.
- **idx_descarga_archivo_id (archivo_id)**: Trackear uso de un archivo.
- **idx_descarga_status (status)**: Monitoreo y filtrado por estado.
- **idx_descarga_user_status (user_id, status)**: Dashboard de usuario con filtros de estado.

**Consultas optimizadas**:
- Historial del usuario: `SELECT * FROM "Descarga" WHERE user_id = $1 ORDER BY created_at DESC`
- Descargas activas: `SELECT * FROM "Descarga" WHERE status = 'processing'`
- Descargas por usuario y estado: `SELECT * FROM "Descarga" WHERE user_id = $1 AND status = $2`

### Embedding

- **idx_embedding_archivo_id (archivo_id)**: Obtener embeddings de un archivo.
- **idx_embedding_model (model_name)**: Filtrar por modelo de embedding.
- **idx_embedding_hnsw (embedding vector_cosine_ops)**: Busqueda de similitud semantica con HNSW.

**Consultas optimizadas**:
- Embeddings de un archivo: `SELECT * FROM "Embedding" WHERE archivo_id = $1`
- Busqueda semantica: `SELECT * FROM "Embedding" ORDER BY embedding <=> $1 LIMIT 10`

### WorkerExecution

- **idx_workerexecution_download_id (download_id)**: Historial de ejecuciones de una descarga.
- **idx_workerexecution_status (status)**: Monitoreo de workers activos.
- **idx_workerexecution_worker (worker_name)**: Metricas por tipo de worker.

**Consultas optimizadas**:
- Ejecuciones de una descarga: `SELECT * FROM "WorkerExecution" WHERE download_id = $1 ORDER BY attempt`
- Workers activos: `SELECT * FROM "WorkerExecution" WHERE status = 'running'`
- Metricas por worker: `SELECT worker_name, AVG(duration_ms) FROM "WorkerExecution" GROUP BY worker_name`

### AuditEvent

- **idx_auditevent_download_id (download_id)**: Auditoria de una descarga especifica.
- **idx_auditevent_event_type (event_type)**: Filtrado por tipo de evento.
- **idx_auditevent_worker_execution (worker_execution_id)**: Eventos de una ejecucion.
- **idx_auditevent_occurred_at (occurred_at)**: Rangos de tiempo y ordenamiento.
- **idx_auditevent_payload (payload GIN)**: Busqueda dentro del JSONB.

**Consultas optimizadas**:
- Auditoria de descarga: `SELECT * FROM "AuditEvent" WHERE download_id = $1 ORDER BY occurred_at`
- Eventos por tipo: `SELECT * FROM "AuditEvent" WHERE event_type = 'DownloadCompleted'`
- Busqueda en payload: `SELECT * FROM "AuditEvent" WHERE payload @> '{"key": "value"}'`

## Particionamiento

### AuditEvent

- **Estrategia**: Particionado por rango mensual sobre `occurred_at`
- **Razon**: Millones de eventos de auditoria acumulados. Sin particionamiento, cada consulta escanea toda la tabla.
- **Beneficios**:
  - Consultas recientes solo tocan la particion activa
  - Archivado a Iceberg/Parquet por particion antigua
  - DROP de particiones antiguas es instantaneo y libera espacio
- **Mantenimiento**: Crear nuevas particiones mensualmente con un cron job

## Recomendaciones de VACUUM y ANALYZE

### VACUUM

- **Tablas de alta escritura**: `Descarga`, `WorkerExecution`, `AuditEvent`
- **Frecuencia recomendada**: 
  - `VACUUM ANALYZE` diario en horas de baja carga
  - `VACUUM FULL` mensual durante mantenimiento programado
- **Configuracion de autovacuum**: Asegurar que autovacuum este habilitado y ajustar thresholds para tablas grandes

### ANALYZE

- **Frecuencia**: Despues de cargas masivas (bulk inserts) o al menos una vez al dia
- **Tablas criticas**: Todas, especialmente aquellas con distribucion de datos cambiante
- **Comando**: `ANALYZE "Descarga"; ANALYZE "AuditEvent";` etc.

## Recomendaciones para pgvector

1. **Dimension del vector**: El modelo usa 768 dimensiones. Si se cambia de modelo, ajustar `VECTOR(N)`.
2. **HNSW vs IVFFlat**: 
   - HNSW elegido por mayor precision y mejor rendimiento en busqueda
   - Mayor costo de construccion aceptable dado que inserciones son menos frecuentes que consultas
3. **Mantenimiento del indice HNSW**:
   - No requiere reconstruccion manual como IVFFlat
   - Monitorear uso de memoria en servidores con recursos limitados
4. **Busqueda eficiente**:
   - Siempre incluir un `LIMIT` en consultas de similitud
   - Considerar pre-filtrado por `model_name` o `archivo_id` antes de busqueda vectorial

## Posibles Mejoras Futuras

### Sharding/Particionamiento adicional

- **Descarga**: Considerar particionado por `user_id` o `created_at` si crece masivamente
- **WorkerExecution**: Particionado por `created_at` o `download_id`

### Indices adicionales

- **Archivo**: Indice parcial para archivos con `reference_count > 0` (candidatos a limpieza)
- **Descarga**: Indice parcial para `status = 'processing'`, `status = 'pending'` (monitoreo activo)
- **Usuario**: Indice para monitoreo de cuota `storage_used_bytes > storage_quota_bytes`

### Materialized Views

- Dashboard de metricas: Crear vistas materializadas para reportes frecuentes (descargas por dia, usuarios activos, etc.)
- Actualizacion periodica con REFRESH MATERIALIZED VIEW CONCURRENTLY

### Archivado

- Implementar un proceso automatizado (cron job o worker) que:
  1. Identifique particiones de `AuditEvent` mayores a 90 dias
  2. Exporte datos a Parquet en Iceberg
  3. Ejecute `DROP PARTITION` para liberar espacio en OLTP
  4. Mantenga un catalogo de datos archivados para consultas historicas via DuckDB

## Conclusion

El modelo esta disenado para escalar hasta millones de registros con:
- Indices estrategicos para patrones de consulta conocidos
- Particionamiento agresivo de tablas de alta volumen (AuditEvent)
- Triggers para mantener integridad sin sacrificar performance
- TIMESTAMPTZ para consistencia en sistemas distribuidos
- UUID v7 para localidad de indices y trazabilidad

El mantenimiento regular (VACUUM, ANALYZE, creacion de particiones) sera la clave para mantener el rendimiento a largo plazo.
