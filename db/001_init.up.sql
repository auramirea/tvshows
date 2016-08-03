CREATE TABLE tvshow (
  id serial NOT NULL,
  name        VARCHAR(255) NOT NULL,
  url         VARCHAR(255) NOT NULL,
  image       VARCHAR(255),
  rating      FLOAT,
  created_at  timestamp with time zone,
  updated_at  timestamp with time zone,
  deleted_at  timestamp with time zone,
  CONSTRAINT tvshow_pkey PRIMARY KEY (id)
);

CREATE TABLE appuser
(
  id          serial   NOT NULL,
  created_at  timestamp with time zone,
  updated_at  timestamp with time zone,
  deleted_at  timestamp with time zone,
  first_name  character varying(255),
  last_name   character varying(255),
  email       character varying(255),
  CONSTRAINT appuser_email_unique UNIQUE (email),
  CONSTRAINT appuser_pkey PRIMARY KEY (id)
)