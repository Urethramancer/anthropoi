package anthropoi

const groupTables = `BEGIN WORK;
CREATE TABLE public.groups
(
	-- id auto-increments
	id serial NOT NULL,
	-- name of the group.
	name character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	created timestamp with time zone,
	CONSTRAINT group_id_pkey PRIMARY KEY (id),
	CONSTRAINT group_name_unique UNIQUE (name)
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_groups_timestamp ON public.groups;
CREATE TRIGGER trigger_groups_timestamp
	BEFORE INSERT ON public.groups
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- Permissions for all groups
CREATE TABLE public.permissions
(
	id serial NOT NULL,
	groupid integer NOT NULL,
	name character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	CONSTRAINT permissions_id_pkey PRIMARY KEY (id),
	CONSTRAINT permissions_groupid_fkey FOREIGN KEY (groupid)
	REFERENCES public.groups (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Users can have multiple membership relations per domain to this table,
-- since groups can have multiple application-defined permissions.
CREATE TABLE public.membership
(
	-- user this membership is for.
	userid integer NOT NULL,
	-- profile this membership is for.
	profileid integer NOT NULL,
	-- permission for this user:profile:group combination
	permission integer NOT NULL,
	-- memberships need a creation time too
	created timestamp with time zone,
	-- membership:users relationship
	CONSTRAINT membership_user_fkey FOREIGN KEY (userid)
	REFERENCES public.users (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	-- membership:profile relationship
	CONSTRAINT membership_profile_fkey FOREIGN KEY (profileid)
	REFERENCES public.profiles (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	-- membership:permissions relationship
	CONSTRAINT membership_permissions_fkey FOREIGN KEY (permission)
	REFERENCES public.permissions (id) MATCH SIMPLE
		ON UPDATE CASCADE
		ON DELETE CASCADE
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_membership_timestamp ON public.membership;
CREATE TRIGGER trigger_membership_timestamp
	BEFORE INSERT ON public.membership
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
COMMIT WORK;
`
