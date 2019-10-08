// Package inspect contains functions to query table structures, constraints and keys.
package inspect

// Gets the columns, types and defaults for a table.
const tableQuery = `select column_name, data_type, character_maximum_length, column_default, is_nullable, numeric_precision, table_schema
from information_schema.columns where table_name='{TABLE}';`

// Gets the primary key for a table.
const pkQuery = `SELECT c.column_name AS col
FROM information_schema.key_column_usage AS c
LEFT JOIN information_schema.table_constraints AS t
ON t.constraint_name = c.constraint_name
WHERE t.table_name = 'users' AND t.constraint_type = 'PRIMARY KEY';
`

// Gets the unique columns for a table.
const uniqueQuery = `SELECT c.column_name AS col
FROM information_schema.key_column_usage AS c
LEFT JOIN information_schema.table_constraints AS t
ON t.constraint_name = c.constraint_name
WHERE t.table_name = 'users' AND t.constraint_type = 'UNIQUE';
`

// Gets the foreign key(s) for a table.
const fkQuery = `select kcu.table_schema || '.' || kcu.table_name as table,
kcu.column_name as key,
rel_kcu.table_schema || '.' || rel_kcu.table_name as dest_table,
rel_kcu.column_name as dest_col,
kcu.constraint_name,
-- NONE = MATCH SIMPLE
match_option,
-- ON UPDATE behaviour (CASCADE etc.)
update_rule,
-- ON DELETE behaviour
delete_rule
from information_schema.table_constraints tco
join information_schema.key_column_usage kcu
   on tco.constraint_schema=kcu.constraint_schema
   and tco.constraint_name=kcu.constraint_name
join information_schema.referential_constraints rco
   on tco.constraint_schema=rco.constraint_schema
   and tco.constraint_name=rco.constraint_name
join information_schema.key_column_usage rel_kcu
   on rco.unique_constraint_schema=rel_kcu.constraint_schema
   and rco.unique_constraint_name=rel_kcu.constraint_name
   and kcu.ordinal_position=rel_kcu.ordinal_position
where tco.constraint_type='FOREIGN KEY'
and kcu.table_name='{TABLE}'
order by kcu.table_schema,
  kcu.table_name;
`
