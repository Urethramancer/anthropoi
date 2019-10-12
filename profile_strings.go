package anthropoi

const profileTable = `BEGIN WORK;
CREATE TABLE public.profiles
(
	-- id auto-increments
	id serial NOT NULL,
	userid integer NOT NULL,
	domain character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	created timestamp with time zone,
	data json NOT NULL DEFAULT '{}'::json,
	CONSTRAINT con_profiles_pkey PRIMARY KEY (id),
	CONSTRAINT key_profiles_userid_fkey FOREIGN KEY (userid)
	REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_profiles_timestamp ON public.profiles;
CREATE TRIGGER trigger_profiles_timestamp
	BEFORE INSERT ON public.profiles
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
COMMIT WORK;
`
