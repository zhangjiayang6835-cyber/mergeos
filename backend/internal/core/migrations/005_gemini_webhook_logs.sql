CREATE TABLE IF NOT EXISTS gemini_webhook_logs (
  id text PRIMARY KEY,
  delivery_id text NOT NULL DEFAULT '',
  event_name text NOT NULL DEFAULT '',
  action text NOT NULL DEFAULT '',
  repository text NOT NULL DEFAULT '',
  pull_number integer NOT NULL DEFAULT 0,
  sender text NOT NULL DEFAULT '',
  status text NOT NULL DEFAULT '',
  status_code integer NOT NULL DEFAULT 0,
  error text NOT NULL DEFAULT '',
  comment_url text NOT NULL DEFAULT '',
  key_id text NOT NULL DEFAULT '',
  labels jsonb NOT NULL DEFAULT '[]'::jsonb,
  duration_millis bigint NOT NULL DEFAULT 0,
  received_at timestamptz NOT NULL DEFAULT now(),
  completed_at timestamptz
);

CREATE INDEX IF NOT EXISTS gemini_webhook_logs_received_at_idx ON gemini_webhook_logs (received_at DESC);
CREATE INDEX IF NOT EXISTS gemini_webhook_logs_delivery_id_idx ON gemini_webhook_logs (delivery_id);
CREATE INDEX IF NOT EXISTS gemini_webhook_logs_status_idx ON gemini_webhook_logs (status);
