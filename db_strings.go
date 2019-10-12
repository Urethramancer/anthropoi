package anthropoi

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
