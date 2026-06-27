# Modelo de Datos OLTP - Youtube Downloader

## Introducción

Este documento describe el modelo de datos OLTP (Online Transaction Processing) para la plataforma de gestión y descarga inteligente de videos. El sistema está diseñado para soportar millones de registros con alta concurrencia, utilizando arquitectura distribuida.

## Convenciones Generales

- **Claves primarias**: UUID v7 generado en la aplicación Go
- **Timestamps**: Todos los campos de fecha/hora usan `TIMESTAMPTZ`
- **Patrón de auditoría**:
  - `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
  - `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()` (para entidades modificables)
  - `occurred_at TIMESTAMPTZ` (para eventos y auditoría)
- **Estados y tipos**: Utilizan ENUMs cuando el conjunto de valores es estable
- **Integridad referencial**: Triggers para validaciones complejas entre tablas

## Entidades

### 1. Usuario

Representa una cuenta registrada en la plataforma.

**Propósito**: Gestionar la identidad, autenticación, cuotas de almacenamiento y estado de la cuenta del usuario.

**Campos**:
- `id` (UUID): Identificador único generado por la aplicación Go
- `email` (VARCHAR): Correo electrónico del usuario, único
- `username` (VARCHAR): Nombre de usuario, único
- `password_hash` (VARCHAR): Hash de la contraseña
- `display_name` (VARCHAR): Nombre visible del usuario
- `plan` (ENUM): Plan de suscripción (free, basic, premium, enterprise)
- `status` (ENUM): Estado de la cuenta (active, suspended, cancelled)
- `storage_used_bytes` (BIGINT): Bytes de almacenamiento utilizados
- `storage_quota_bytes` (BIGINT): Cuota de almacenamiento asignada
- `created_at` (TIMESTAMPTZ): Fecha de creación
- `updated_at` (TIMESTAMPTZ): Fecha de última modificación
- `last_login_at` (TIMESTAMPTZ): Fecha del último inicio de sesión

**Decisiones de diseño**:
- `password_hash` almacena el hash, nunca la contraseña en texto plano
- `storage_used_bytes` y `storage_quota_bytes` en BIGINT para soportar rangos amplios
- `plan` y `status` como ENUM para integridad de datos y reducción de errores
- `last_login_at` actualizado por la aplicación, sin trigger automático

---

### 2. Video

Representa el contenido lógico de un video de YouTube.

**Propósito**: Almacenar metadatos del video, independiente de cualquier archivo físico o descarga.

**Campos**:
- `id` (UUID): Identificador único
- `youtube_video_id` (VARCHAR): ID de YouTube del video (ej: 'dQw4w9WgXcQ')
- `title` (VARCHAR): Título del video
- `channel_name` (VARCHAR): Nombre del canal de YouTube
- `duration_seconds` (INTEGER): Duración en segundos
- `published_at` (TIMESTAMPTZ): Fecha de publicación en YouTube
- `created_at` (TIMESTAMPTZ): Fecha de creación en el sistema

**Decisiones de diseño**:
- No contiene información técnica de archivos (codec, bitrate, etc.) — eso pertenece a `Archivo`
- `youtube_video_id` es único para evitar duplicación de videos
- `published_at` usa TIMESTAMPTZ para consistencia con fechas de YouTube
- Sin `updated_at` porque los metadatos del video no se actualizan una vez creados

---

### 3. Archivo

Representa una representación física de un Video.

**Propósito**: Gestionar los archivos físicos almacenados en Cloudflare R2, con deduplicación mediante `reference_count`.

**Campos**:
- `id` (UUID): Identificador único
- `video_id` (UUID): Referencia al Video al que pertenece
- `sha256` (VARCHAR): Hash SHA256 del archivo, único para deduplicación
- `object_key` (VARCHAR): Clave S3 del objeto en Cloudflare R2
- `storage_provider` (ENUM): Proveedor de almacenamiento (r2, s3, etc.)
- `mime_type` (VARCHAR): Tipo MIME del archivo
- `file_type` (ENUM): Tipo de archivo (video, audio, subtitle, thumbnail)
- `codec` (VARCHAR): Códec del archivo
- `width` (INTEGER): Ancho en píxeles (si aplica)
- `height` (INTEGER): Alto en píxeles (si aplica)
- `fps` (NUMERIC): Frames por segundo (si aplica)
- `bitrate` (INTEGER):Bitrate en bps (si aplica)
- `size_bytes` (BIGINT): Tamaño en bytes
- `reference_count` (INTEGER): Contador de referencias para deduplicación
- `created_at` (TIMESTAMPTZ): Fecha de creación

**Decisiones de diseño**:
- `reference_count` mantiene la cuenta de descargas que usan este archivo, operacional aunque no analítico
- `sha256` único permite deduplicación a nivel de contenido
- Sin `updated_at` porque los archivos físicos son inmutables una vez creados
- Campos de resolución (width, height) son NULL cuando no aplican (ej: audio)

---

### 4. Descarga

Representa la solicitud de descarga realizada por un usuario.

**Propósito**: Rastrear qué usuario solicitó qué video, en qué formato/quality, y cuál es el estado del proceso.

**Campos**:
- `id` (UUID): Identificador único
- `user_id` (UUID): Usuario que solicitó la descarga
- `video_id` (UUID): Video solicitado (contenido lógico)
- `archivo_id` (UUID): Archivo físico entregado
- `status` (ENUM): Estado de la descarga (pending, processing, completed, failed, cancelled)
- `requested_quality` (VARCHAR): Calidad solicitada (ej: '1080p', '720p')
- `requested_format` (VARCHAR): Formato solicitado (ej: 'mp4', 'mp3')
- `started_at` (TIMESTAMPTZ): Fecha de inicio del procesamiento
- `finished_at` (TIMESTAMPTZ): Fecha de finalización
- `created_at` (TIMESTAMPTZ): Fecha de creación de la solicitud
- `updated_at` (TIMESTAMPTZ): Fecha de última modificación

**Decisiones de diseño**:
- `video_id` y `archivo_id` son ambos obligatorios; un trigger garantiza que el archivo pertenezca al video
- `archivo_id` puede ser NULL inicialmente y se actualiza cuando el archivo es seleccionado/generado
- `started_at` y `finished_at` permiten calcular la duración del procesamiento
- `updated_at` es actualizado automáticamente por trigger

---

### 5. Embedding

Cada Archivo puede tener múltiples embeddings generados por modelos de IA.

**Propósito**: Almacenar representaciones vectoriales de archivos para búsqueda semántica con pgvector.

**Campos**:
- `id` (UUID): Identificador único
- `archivo_id` (UUID): Archivo al que pertenece el embedding
- `model_name` (VARCHAR): Nombre del modelo usado para generar el embedding
- `embedding` (VECTOR): Vector de embeddings de pgvector
- `created_at` (TIMESTAMPTZ): Fecha de creación

**Decisiones de diseño**:
- `embedding` usa el tipo `VECTOR` de pgvector
- Índice HNSW para búsqueda de similitud eficiente
- Sin `updated_at` porque los embeddings se regeneran completamente, no se actualizan
- Cada archivo puede tener múltiples embeddings de diferentes modelos o propósitos (OCR, speech, descripción)

---

### 6. WorkerExecution

Registra cada ejecución de un worker del pipeline de procesamiento.

**Propósito**: Trazabilidad completa del ciclo de vida de cada intento de procesamiento.

**Campos**:
- `id` (UUID): Identificador único
- `download_id` (UUID): Descarga asociada
- `worker_name` (VARCHAR): Nombre del worker ejecutado
- `status` (ENUM): Estado de la ejecución (pending, running, completed, failed, retried)
- `attempt` (INTEGER): Número de intento
- `started_at` (TIMESTAMPTZ): Fecha de inicio
- `finished_at` (TIMESTAMPTZ): Fecha de finalización
- `duration_ms` (INTEGER): Duración en milisegundos
- `error_message` (TEXT): Mensaje de error (si aplica)

**Decisiones de diseño**:
- Cada intento de ejecución crea una nueva fila (no se actualizan registros)
- `attempt` numera cada retry para la misma Descarga
- `duration_ms` es redundante con `started_at` y `finished_at` pero optimiza lecturas frecuentes
- Sin `updated_at` porque las ejecuciones son inmutables una vez finalizadas

---

### 7. AuditEvent

Tabla de auditoría que registra todas las transiciones del pipeline.

**Propósito**: Auditoría completa e inmutable de todos los eventos del sistema.

**Campos**:
- `id` (UUID): Identificador único
- `download_id` (UUID): Descarga relacionada
- `worker_execution_id` (UUID): Ejecución de worker relacionada
- `event_type` (ENUM): Tipo de evento (DownloadRequested, DownloadStarted, DownloadCompleted, MetadataStarted, MetadataCompleted, OCRStarted, OCRCompleted, EmbeddingStarted, EmbeddingCompleted, OptimizationStarted, OptimizationCompleted, Stored, Completed, Failed, Retry)
- `status` (ENUM): Estado del evento (pending, processing, completed, failed)
- `payload` (JSONB): Datos adicionales del evento
- `error_message` (TEXT): Mensaje de error
- `occurred_at` (TIMESTAMPTZ): Fecha de ocurrencia del evento

**Decisiones de diseño**:
- Particionado por rango mensual en `occurred_at` para escalabilidad
- `payload` JSONB permite flexibilidad sin modificar el esquema
- Los eventos son inmutables; no se actualizan ni eliminan
- `occurred_at` (no `created_at`) para consistencia con semántica de eventos
- Ventana operacional de 90 días en OLTP; histórico completo en Lakehouse

---

## Relaciones

```
Usuario 1 ---- N Descarga
Video   1 ---- N Archivo
Video   1 ---- N Descarga
Archivo 1 ---- N Descarga
Archivo 1 ---- N Embedding
Descarga 1 ---- N WorkerExecution
WorkerExecution 1 ---- N AuditEvent
```

### Cardinalidades

- **Usuario - Descarga**: Un usuario puede tener múltiples descargas. Cada descarga pertenece a un único usuario.
- **Video - Archivo**: Un video puede tener múltiples archivos (diferentes calidades/formatos). Cada archivo pertenece a un único video.
- **Video - Descarga**: Un video puede ser descargado múltiples veces. Cada descarga solicita un único video.
- **Archivo - Descarga**: Un archivo puede ser referenciado por múltiples descargas. Cada descarga referencia un único archivo final.
- **Archivo - Embedding**: Un archivo puede tener múltiples embeddings. Cada embedding pertenece a un único archivo.
- **Descarga - WorkerExecution**: Una descarga puede tener múltiples ejecuciones de workers. Cada ejecución pertenece a una única descarga.
- **WorkerExecution - AuditEvent**: Una ejecución de worker puede generar múltiples eventos de auditoría. Cada evento de auditoría puede estar asociado a una única ejecución.

## Decisiones Arquitectónicas

1. **UUID v7**: Generado en la aplicación Go para mejor localidad en índices B-Tree y compatibilidad con sistemas distribuidos
2. **TIMESTAMPTZ**: Siempre usado para evitar ambigüedades de zonas horarias en sistema distribinerar
3. **ENUMs**: Usados para estados y tipos con conjuntos estables de valores, mejorando integridad y auto-documentación
4. **Triggers de integridad cross-table**: Implementados donde CHECK constraints no son suficientes (ej: descarga.video_id = archivo.video_id)
5. **Patrón de auditoría genérico**: `created_at` + `updated_at` + trigger automático para coherencia en todas las tablas que lo requieren
6. **Particionado de AuditEvent**: Mensual por `occurred_at` para mantener el OLTP operativo con millones de registros
