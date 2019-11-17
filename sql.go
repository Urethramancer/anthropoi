package anthropoi

const DefaultName = "accounts"
const databaseDefinitions = `CREATE DATABASE {NAME};`
const databaseTriggers = `BEGIN WORK;
	-- We'll trigger creation timestamp setting in a few places.
	CREATE OR REPLACE FUNCTION trigger_set_timestamp()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.created = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	COMMIT WORK;
`

const varTables = `CREATE TABLE public.variables
(
	key character varying(100) COLLATE pg_catalog."default" NOT NULL,
	value text COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	CONSTRAINT var_prim PRIMARY KEY (key),
	CONSTRAINT var_key_is_unique UNIQUE (key)
) WITH (OIDS = FALSE) TABLESPACE pg_default;
`

const flagTables = `CREATE TABLE public.flags
(
	key character varying(100) COLLATE pg_catalog."default" NOT NULL,
	flag boolean NOT NULL DEFAULT false,
	CONSTRAINT flags_prim PRIMARY KEY (key),
	CONSTRAINT flags_key_is_unique UNIQUE (key)
) WITH (OIDS = FALSE) TABLESPACE pg_default;
`

const groupTables = `BEGIN WORK;
CREATE TABLE IF NOT EXISTS public.sites
(
	-- id auto-increments
	id serial NOT NULL,
	name character varying(100) COLLATE pg_catalog."default" NOT NULL,
	created timestamp with time zone,
	CONSTRAINT sites_id_pkey PRIMARY KEY (id),
	CONSTRAINT sites_name_unique UNIQUE (name)
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_sites_timestamp ON public.sites;
CREATE TRIGGER trigger_sites_timestamp
	BEFORE INSERT ON public.sites
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- User to site connections
CREATE TABLE IF NOT EXISTS public.membership
(
	-- user this membership is for.
	userid integer NOT NULL,
	-- site the user is a member of
	siteid integer NOT NULL,
	-- how long a user has been a member
	created timestamp with time zone,
	-- primary key
	CONSTRAINT member_combined_pkey PRIMARY KEY(userid,siteid),
	-- membership:users relationship
	CONSTRAINT member_users_fkey FOREIGN KEY (userid)
	REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT member_sites_fkey FOREIGN KEY (siteid)
	REFERENCES public.sites (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;
	
-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_member_timestamp ON public.membership;
CREATE TRIGGER trigger_member_timestamp
	BEFORE INSERT ON public.membership
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS public.groups
(
	-- id auto-increments
	id serial NOT NULL,
	-- name of the group (not unique, since multiple sites could have the same groups)
	name character varying(100) COLLATE pg_catalog."default" NOT NULL,
	-- site this group is for
	site integer NOT NULL,
	created timestamp with time zone,
	CONSTRAINT group_id_pkey PRIMARY KEY (id),
	CONSTRAINT groups_site_fkey FOREIGN KEY (site)
	REFERENCES public.sites (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_groups_timestamp ON public.groups;
CREATE TRIGGER trigger_groups_timestamp
	BEFORE INSERT ON public.groups
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- User to group connections
CREATE TABLE IF NOT EXISTS public.roles
(
	userid integer NOT NULL,
	groupid integer NOT NULL,
	CONSTRAINT roles_combined_pkey PRIMARY KEY (userid,groupid),
	CONSTRAINT roles_userid_fkey FOREIGN KEY (userid)
	REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT roles_groupid_fkey FOREIGN KEY (groupid)
	REFERENCES public.groups (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Permissions for all groups
CREATE TABLE IF NOT EXISTS public.permissions
(
	id serial NOT NULL,
	groupid integer NOT NULL,
	name character varying(100) COLLATE pg_catalog."default" NOT NULL,
	CONSTRAINT permissions_id_pkey PRIMARY KEY (id,groupid),
	CONSTRAINT permissions_groupid_fkey FOREIGN KEY (groupid)
	REFERENCES public.groups (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT permissions_name_unique UNIQUE (name)
) WITH (OIDS = FALSE) TABLESPACE pg_default;

COMMIT WORK;
`

const userTable = `BEGIN WORK;
CREATE TABLE IF NOT EXISTS public.users
(
	-- id auto-increments
	id serial NOT NULL,
	-- username for logins across any site in the system.
	username character varying(200) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- password is the hash from bcrypt. 60 is supposed to be a good length for the next millennium.
	password character varying(120) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- salt is protection against rainbow tables.
	salt character varying(32) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Email is required for verification and resetting.
	email character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	created timestamp with time zone,
	locked boolean NOT NULL DEFAULT false,
	admin boolean NOT NULL DEFAULT false,
	first character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	last character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	data json NOT NULL DEFAULT '{}'::json,
	tokens json NOT NULL DEFAULT '{}'::json,
	CONSTRAINT key_users_pkey PRIMARY KEY (id),
	CONSTRAINT text_username_unique UNIQUE (username)
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_users_timestamp ON public.users;
CREATE TRIGGER trigger_users_timestamp
	BEFORE INSERT ON public.users
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS public.resetkeys
(
	key character varying(64) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	account integer NOT NULL,
	CONSTRAINT resetkeys_pkey PRIMARY KEY(key),
	CONSTRAINT resetkeys_account_fkey FOREIGN KEY (account)
	REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	expiry timestamp with time zone NOT NULL
) WITH (OIDS = FALSE) TABLESPACE pg_default;
COMMIT WORK;
`

const aliasesTable = `CREATE TABLE public.aliases
(
	alias character varying(200) COLLATE pg_catalog."default" NOT NULL,
	target character varying(200) COLLATE pg_catalog."default" NOT NULL,
	CONSTRAINT alias_prim PRIMARY KEY (alias, target),
	CONSTRAINT alias_is_unique UNIQUE (alias)
)

`
