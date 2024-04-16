CREATE TABLE IF NOT EXISTS public."user" (
	id varchar NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz NULL,
	provider text NOT NULL,
	handle text NOT NULL,
	"password" text NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.message (
	id text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz NULL,
	scheduled_at timestamptz NULL,
	sent_at timestamptz NULL,
	"type" text NOT NULL,
	criteria text NULL,
	"content" text NOT NULL,
	CONSTRAINT message_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS TABLE public.user_message (
	id text NOT NULL,
	user_id text NOT NULL,
	message_id text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz NULL,
	read_at timestamptz NULL,
	CONSTRAINT user_message_pkey PRIMARY KEY (id)
);
