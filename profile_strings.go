package anthropoi

const profileTable = `CREATE TABLE public.profiles
(
	-- id auto-increments
	id serial NOT NULL,
	userid integer NOT NULL,
	domain character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	created timestamp with time zone,
	data json NOT NULL DEFAULT '{}'::json,
	CONSTRAINT profiles_pkey PRIMARY KEY (id),
	CONSTRAINT profiles_user_fkey FOREIGN KEY (userid)
	REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
CREATE TRIGGER set_profiles_timestamp
	BEFORE INSERT ON public.profiles
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
`
