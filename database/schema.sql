--
-- Youtube Downloader - Schema PostgreSQL para Neon
-- Generado para produccion con pgvector, particionamiento y triggers
--

-- =========================================
-- EXTENSIONES
-- =========================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "vector";

-- =========================================
-- TIPOS ENUM
-- =========================================

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'plan_type') THEN
        CREATE TYPE plan_type AS ENUM ('free', 'basic', 'premium', 'enterprise');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('active', 'suspended', 'cancelled');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'storage_provider') THEN
        CREATE TYPE storage_provider AS ENUM ('r2', 's3', 'gcs');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'file_type') THEN
        CREATE TYPE file_type AS ENUM ('video', 'audio', 'subtitle', 'thumbnail');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'download_status') THEN
        CREATE TYPE download_status AS ENUM ('pending', 'processing', 'completed', 'failed', 'cancelled');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'worker_status') THEN
        CREATE TYPE worker_status AS ENUM ('pending', 'running', 'completed', 'failed', 'retried');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
        CREATE TYPE event_type AS ENUM (
            'DownloadRequested', 'DownloadStarted', 'DownloadCompleted',
            'MetadataStarted', 'MetadataCompleted',
            'OCRStarted', 'OCRCompleted',
            'EmbeddingStarted', 'EmbeddingCompleted',
            'OptimizationStarted', 'OptimizationCompleted',
            'Stored', 'Completed', 'Failed', 'Retry'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_status') THEN
        CREATE TYPE event_status AS ENUM ('pending', 'processing', 'completed', 'failed');
    END IF;
END$$;

-- =========================================
-- SECUENCIAS PARA ATTEMPT EN WorkerExecution
-- =========================================

CREATE SEQUENCE IF NOT EXISTS worker_execution_attempt_seq START 1;

-- =========================================
-- TABLAS
-- =========================================

-- 1. Usuario
CREATE TABLE IF NOT EXISTS "Usuario" (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    plan plan_type NOT NULL DEFAULT 'free',
    status user_status NOT NULL DEFAULT 'active',
    storage_used_bytes BIGINT NOT NULL DEFAULT 0,
    storage_quota_bytes BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at TIMESTAMPTZ
);

COMMENT ON TABLE "Usuario" IS 'Cuentas registradas en la plataforma';
COMMENT ON COLUMN "Usuario".id IS 'Identificador unico (UUID v7 generado en Go)';
COMMENT ON COLUMN "Usuario".email IS 'Correo electronico del usuario';
COMMENT ON COLUMN "Usuario".username IS 'Nombre de usuario unico';
COMMENT ON COLUMN "Usuario".password_hash IS 'Hash de la contrasena (nunca texto plano)';
COMMENT ON COLUMN "Usuario".display_name IS 'Nombre visible del usuario';
COMMENT ON COLUMN "Usuario".plan IS 'Plan de suscripcion';
COMMENT ON COLUMN "Usuario".status IS 'Estado de la cuenta';
COMMENT ON COLUMN "Usuario".storage_used_bytes IS 'Bytes de almacenamiento utilizados';
COMMENT ON COLUMN "Usuario".storage_quota_bytes IS 'Cuota de almacenamiento asignada';
COMMENT ON COLUMN "Usuario".created_at IS 'Fecha de creacion';
COMMENT ON COLUMN "Usuario".updated_at IS 'Fecha de ultima modificacion';
COMMENT ON COLUMN "Usuario".last_login_at IS 'Fecha del ultimo inicio de sesion';

-- =========================================

-- 2. Video
CREATE TABLE IF NOT EXISTS "Video" (
    id UUID PRIMARY KEY,
    youtube_video_id VARCHAR(20) NOT NULL UNIQUE,
    title VARCHAR(500) NOT NULL,
    channel_name VARCHAR(255),
    duration_seconds INTEGER,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON TABLE "Video" IS 'Contenido logico de videos de YouTube';
COMMENT ON COLUMN "Video".id IS 'Identificador unico';
COMMENT ON COLUMN "Video".youtube_video_id IS 'ID de YouTube del video';
COMMENT ON COLUMN "Video".title IS 'Titulo del video';
COMMENT ON COLUMN "Video".channel_name IS 'Nombre del canal de YouTube';
COMMENT ON COLUMN "Video".duration_seconds IS 'Duracion en segundos';
COMMENT ON COLUMN "Video".published_at IS 'Fecha de publicacion en YouTube';
COMMENT ON COLUMN "Video".created_at IS 'Fecha de creacion en el sistema';

-- =========================================

-- 3. Archivo
CREATE TABLE IF NOT EXISTS "Archivo" (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL,
    sha256 VARCHAR(64) NOT NULL UNIQUE,
    object_key VARCHAR(500) NOT NULL,
    storage_provider storage_provider NOT NULL DEFAULT 'r2',
    mime_type VARCHAR(100),
    file_type file_type NOT NULL,
    codec VARCHAR(50),
    width INTEGER,
    height INTEGER,
    fps NUMERIC(8,2),
    bitrate INTEGER,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    reference_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    CONSTRAINT fk_archivo_video FOREIGN KEY (video_id) REFERENCES "Video"(id) ON DELETE RESTRICT
);

COMMENT ON TABLE "Archivo" IS 'Representacion fisica de un Video';
COMMENT ON COLUMN "Archivo".id IS 'Identificador unico';
COMMENT ON COLUMN "Archivo".video_id IS 'Video al que pertenece el archivo';
COMMENT ON COLUMN "Archivo".sha256 IS 'Hash SHA256 del archivo para deduplicacion';
COMMENT ON COLUMN "Archivo".object_key IS 'Clave S3 del objeto en Cloudflare R2';
COMMENT ON COLUMN "Archivo".storage_provider IS 'Proveedor de almacenamiento';
COMMENT ON COLUMN "Archivo".mime_type IS 'Tipo MIME del archivo';
COMMENT ON COLUMN "Archivo".file_type IS 'Tipo de archivo';
COMMENT ON COLUMN "Archivo".codec IS 'Codec del archivo';
COMMENT ON COLUMN "Archivo".width IS 'Ancho en pixeles';
COMMENT ON COLUMN "Archivo".height IS 'Alto en pixeles';
COMMENT ON COLUMN "Archivo".fps IS 'Frames por segundo';
COMMENT ON COLUMN "Archivo".bitrate IS 'Bitrate en bps';
COMMENT ON COLUMN "Archivo".size_bytes IS 'Tamano en bytes';
COMMENT ON COLUMN "Archivo".reference_count IS 'Contador de referencias para deduplicacion';
COMMENT ON COLUMN "Archivo".created_at IS 'Fecha de creacion';

-- =========================================

-- 4. Descarga
CREATE TABLE IF NOT EXISTS "Descarga" (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    video_id UUID NOT NULL,
    archivo_id UUID,
    status download_status NOT NULL DEFAULT 'pending',
    requested_quality VARCHAR(50),
    requested_format VARCHAR(50),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    CONSTRAINT fk_descarga_usuario FOREIGN KEY (user_id) REFERENCES "Usuario"(id) ON DELETE RESTRICT,
    CONSTRAINT fk_descarga_video FOREIGN KEY (video_id) REFERENCES "Video"(id) ON DELETE RESTRICT,
    CONSTRAINT fk_descarga_archivo FOREIGN KEY (archivo_id) REFERENCES "Archivo"(id) ON DELETE RESTRICT
);

COMMENT ON TABLE "Descarga" IS 'Solicitud de descarga realizada por un usuario';
COMMENT ON COLUMN "Descarga".id IS 'Identificador unico';
COMMENT ON COLUMN "Descarga".user_id IS 'Usuario que solicito la descarga';
COMMENT ON COLUMN "Descarga".video_id IS 'Video solicitado (contenido logico)';
COMMENT ON COLUMN "Descarga".archivo_id IS 'Archivo fisico entregado';
COMMENT ON COLUMN "Descarga".status IS 'Estado de la descarga';
COMMENT ON COLUMN "Descarga".requested_quality IS 'Calidad solicitada';
COMMENT ON COLUMN "Descarga".requested_format IS 'Formato solicitado';
COMMENT ON COLUMN "Descarga".started_at IS 'Fecha de inicio del procesamiento';
COMMENT ON COLUMN "Descarga".finished_at IS 'Fecha de finalizacion';
COMMENT ON COLUMN "Descarga".created_at IS 'Fecha de creacion de la solicitud';
COMMENT ON COLUMN "Descarga".updated_at IS 'Fecha de ultima modificacion';

-- =========================================

-- 5. Embedding
CREATE TABLE IF NOT EXISTS "Embedding" (
    id UUID PRIMARY KEY,
    archivo_id UUID NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    embedding VECTOR(768) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    CONSTRAINT fk_embedding_archivo FOREIGN KEY (archivo_id) REFERENCES "Archivo"(id) ON DELETE CASCADE
);

COMMENT ON TABLE "Embedding" IS 'Embeddings de archivos para busqueda semantica';
COMMENT ON COLUMN "Embedding".id IS 'Identificador unico';
COMMENT ON COLUMN "Embedding".archivo_id IS 'Archivo al que pertenece el embedding';
COMMENT ON COLUMN "Embedding".model_name IS 'Nombre del modelo usado para generar el embedding';
COMMENT ON COLUMN "Embedding".embedding IS 'Vector de embeddings (pgvector)';
COMMENT ON COLUMN "Embedding".created_at IS 'Fecha de creacion';

-- =========================================

-- 6. WorkerExecution
CREATE TABLE IF NOT EXISTS "WorkerExecution" (
    id UUID PRIMARY KEY,
    download_id UUID NOT NULL,
    worker_name VARCHAR(255) NOT NULL,
    status worker_status NOT NULL DEFAULT 'pending',
    attempt INTEGER NOT NULL DEFAULT nextval('worker_execution_attempt_seq'),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    duration_ms INTEGER,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    CONSTRAINT fk_workerexecution_descarga FOREIGN KEY (download_id) REFERENCES "Descarga"(id) ON DELETE CASCADE
);

COMMENT ON TABLE "WorkerExecution" IS 'Registro de cada ejecucion de un worker del pipeline';
COMMENT ON COLUMN "WorkerExecution".id IS 'Identificador unico';
COMMENT ON COLUMN "WorkerExecution".download_id IS 'Descarga asociada';
COMMENT ON COLUMN "WorkerExecution".worker_name IS 'Nombre del worker ejecutado';
COMMENT ON COLUMN "WorkerExecution".status IS 'Estado de la ejecucion';
COMMENT ON COLUMN "WorkerExecution".attempt IS 'Numero de intento';
COMMENT ON COLUMN "WorkerExecution".started_at IS 'Fecha de inicio';
COMMENT ON COLUMN "WorkerExecution".finished_at IS 'Fecha de finalizacion';
COMMENT ON COLUMN "WorkerExecution".duration_ms IS 'Duracion en milisegundos';
COMMENT ON COLUMN "WorkerExecution".error_message IS 'Mensaje de error';

-- =========================================

-- 7. AuditEvent (Particionada por rango mensual)
CREATE TABLE IF NOT EXISTS "AuditEvent" (
    id UUID NOT NULL,
    download_id UUID NOT NULL,
    worker_execution_id UUID,
    event_type event_type NOT NULL,
    status event_status NOT NULL DEFAULT 'pending',
    payload JSONB,
    error_message TEXT,
    occurred_at TIMESTAMPTZ NOT NULL,
    
    PRIMARY KEY (id, occurred_at)
) PARTITION BY RANGE (occurred_at);

COMMENT ON TABLE "AuditEvent" IS 'Auditoria completa de todos los eventos del sistema';
COMMENT ON COLUMN "AuditEvent".id IS 'Identificador unico';
COMMENT ON COLUMN "AuditEvent".download_id IS 'Descarga relacionada';
COMMENT ON COLUMN "AuditEvent".worker_execution_id IS 'Ejecucion de worker relacionada';
COMMENT ON COLUMN "AuditEvent".event_type IS 'Tipo de evento';
COMMENT ON COLUMN "AuditEvent".status IS 'Estado del evento';
COMMENT ON COLUMN "AuditEvent".payload IS 'Datos adicionales del evento (JSONB)';
COMMENT ON COLUMN "AuditEvent".error_message IS 'Mensaje de error';
COMMENT ON COLUMN "AuditEvent".occurred_at IS 'Fecha de ocurrencia del evento';

-- Crear particiones iniciales (ejemplo para 2026)
CREATE TABLE IF NOT EXISTS "AuditEvent_2026_01" PARTITION OF "AuditEvent"
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE IF NOT EXISTS "AuditEvent_2026_02" PARTITION OF "AuditEvent"
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE IF NOT EXISTS "AuditEvent_2026_03" PARTITION OF "AuditEvent"
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

-- =========================================
-- INDICES
-- =========================================

-- Usuario: busqueda por email y username
CREATE INDEX IF NOT EXISTS idx_usuario_email ON "Usuario"(email);
CREATE INDEX IF NOT EXISTS idx_usuario_username ON "Usuario"(username);

-- Video: busqueda por youtube_video_id
CREATE INDEX IF NOT EXISTS idx_video_youtube_id ON "Video"(youtube_video_id);
CREATE INDEX IF NOT EXISTS idx_video_title ON "Video"(title);

-- Archivo: busqueda por video_id y sha256
CREATE INDEX IF NOT EXISTS idx_archivo_video_id ON "Archivo"(video_id);
CREATE INDEX IF NOT EXISTS idx_archivo_sha256 ON "Archivo"(sha256);
CREATE INDEX IF NOT EXISTS idx_archivo_object_key ON "Archivo"(object_key);

-- Descarga: busqueda por usuario y estado
CREATE INDEX IF NOT EXISTS idx_descarga_user_id ON "Descarga"(user_id);
CREATE INDEX IF NOT EXISTS idx_descarga_video_id ON "Descarga"(video_id);
CREATE INDEX IF NOT EXISTS idx_descarga_archivo_id ON "Descarga"(archivo_id);
CREATE INDEX IF NOT EXISTS idx_descarga_status ON "Descarga"(status);
CREATE INDEX IF NOT EXISTS idx_descarga_user_status ON "Descarga"(user_id, status);

-- Embedding: busqueda por archivo y modelo
CREATE INDEX IF NOT EXISTS idx_embedding_archivo_id ON "Embedding"(archivo_id);
CREATE INDEX IF NOT EXISTS idx_embedding_model ON "Embedding"(model_name);

-- WorkerExecution: busqueda por descarga y estado
CREATE INDEX IF NOT EXISTS idx_workerexecution_download_id ON "WorkerExecution"(download_id);
CREATE INDEX IF NOT EXISTS idx_workerexecution_status ON "WorkerExecution"(status);
CREATE INDEX IF NOT EXISTS idx_workerexecution_worker ON "WorkerExecution"(worker_name);

-- AuditEvent: busqueda por descarga y tipo de evento (en particiones)
CREATE INDEX IF NOT EXISTS idx_auditevent_download_id ON "AuditEvent"(download_id);
CREATE INDEX IF NOT EXISTS idx_auditevent_event_type ON "AuditEvent"(event_type);
CREATE INDEX IF NOT EXISTS idx_auditevent_worker_execution ON "AuditEvent"(worker_execution_id);
CREATE INDEX IF NOT EXISTS idx_auditevent_occurred_at ON "AuditEvent"(occurred_at);

-- Indice GIN para payload JSONB en AuditEvent
CREATE INDEX IF NOT EXISTS idx_auditevent_payload ON "AuditEvent" USING GIN (payload);

-- =========================================
-- INDICE HNSW PARA EMBEDDINGS
-- =========================================

CREATE INDEX IF NOT EXISTS idx_embedding_hnsw 
ON "Embedding" USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

COMMENT ON INDEX idx_embedding_hnsw IS 'Indice HNSW para busqueda de similitud coseno en embeddings';

-- =========================================
-- TRIGGERS
-- =========================================

-- Trigger generico para actualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Aplicar el trigger a las tablas que lo requieren
CREATE TRIGGER trigger_usuario_updated_at BEFORE UPDATE ON "Usuario"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_descarga_updated_at BEFORE UPDATE ON "Descarga"
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =========================================
-- TRIGGER DE INTEGRIDAD DESCARGA(VIDEO_ID, ARCHIVO_ID)
-- =========================================

CREATE OR REPLACE FUNCTION check_descarga_integridad()
RETURNS TRIGGER AS $$
DECLARE
    v_video_id UUID;
BEGIN
    -- Solo validar si archivo_id no es NULL
    IF NEW.archivo_id IS NOT NULL THEN
        SELECT video_id INTO v_video_id
        FROM "Archivo"
        WHERE id = NEW.archivo_id;
        
        IF v_video_id IS NULL THEN
            RAISE EXCEPTION 'El archivo_id % no existe', NEW.archivo_id;
        END IF;
        
        IF v_video_id != NEW.video_id THEN
            RAISE EXCEPTION 'El archivo_id % no pertenece al video_id %', NEW.archivo_id, NEW.video_id;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_descarga_integridad BEFORE INSERT OR UPDATE ON "Descarga"
    FOR EACH ROW EXECUTE FUNCTION check_descarga_integridad();

-- =========================================
-- COMENTARIOS FINALES
-- =========================================

COMMENT ON EXTENSION "uuid-ossp" IS 'Extension para funciones UUID (mantenida por compatibilidad, los UUID v7 se generan en Go)';
COMMENT ON EXTENSION "pgcrypto" IS 'Extension para funciones criptograficas';
COMMENT ON EXTENSION "vector" IS 'Extension pgvector para soporte de embeddings vectoriales';
