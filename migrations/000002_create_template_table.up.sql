CREATE TABLE IF NOT EXISTS templates (
   id UUID NOT NULL PRIMARY KEY,
   template_name VARCHAR(64) NOT NULL DEFAULT '',
   created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMP WITHOUT TIME ZONE
);
