package anthropoi

const functionDefinitions = `-- We'll trigger creation timestamp setting in a few places.
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	NEW.created = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`
