package anthropoi

import (
	"crypto/rand"
	"math/big"
)

const (
	validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!'#$%&/()=?@*^<>-.:,;|[]{}"
)

func genString(size int) string {
	s := make([]byte, size)
	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(validChars))))
		if err != nil {
			return ""
		}
		c := validChars[n.Int64()]
		s[i] = c
	}
	return string(s)
}

const userTable = `BEGIN WORK;
CREATE TABLE IF NOT EXISTS public.users
(
	-- id auto-increments
	id serial NOT NULL,
	-- username for logins across any site in the system.
	username character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- password is the hash from bcrypt. 60 is supposed to be a good length for the next millennium.
	password character varying(60) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- salt is protection against rainbow tables.
	salt character varying(32) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	-- Email is required for verification and resetting.
	email character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	created timestamp with time zone,
	locked boolean NOT NULL DEFAULT false,
	first character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	last character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT '',
	data json NOT NULL DEFAULT '{}'::json,
	tokens json NOT NULL DEFAULT '{}'::json,
	CONSTRAINT key_users_pkey PRIMARY KEY (id),
	CONSTRAINT text_username_unique UNIQUE (username),
	CONSTRAINT text_email_unique UNIQUE (email)
) WITH (OIDS = FALSE) TABLESPACE pg_default;

-- Set the current timestamp whenever a row is inserted.
DROP TRIGGER IF EXISTS trigger_users_timestamp ON public.users;
CREATE TRIGGER trigger_users_timestamp
	BEFORE INSERT ON public.users
	FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
COMMIT WORK;
`
