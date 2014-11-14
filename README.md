cqlmm
=====

CQL Migrations Manager

This project aims to be a database migrations manager, similar in spirit,
to [goose](https://bitbucket.org/liamstask/goose), but for Cassandra.

Now, unlike Goose (or many command-line tools, for that matter), cqlmm writes
as little output as possible, however this may change in the future, and will
probably be implemented (if at all) with the use of a `-v|-verbose` flag.
Currently, these are the only times you will see cqlmm produce output:

* An error was encountered; in this case, the output will be written to STDERR
* `goose create` will print the path to the newly-created migration to STDOUT
  so you can use it like so:
  
	  $ vim `cqlmm create MyNewMigration cql`
	  


Initializing your migration directory
-------------------------------------

To initialize your migrations environment, all you need to do is run:

	$ cqlmm init
	
cqlmm will create the following directories and files:

* `<basedir>` is the base directory for everything
* `<basedir>/cqlmm.json` the JSON configuration for cqlmm
* `<basedir>/migrations` is the directory where the migrations will be kept

By default, `<basedir>` is set to `db`. If you wish to change this, you can
specify the `-c` flag:

	$ cqlmm -c=my-migrations-directory init
	
Creating a new migration
------------------------

After you have run `cqlmm init`, you can create a new migration by running

	$ cqlmm create NAME FORMAT

This will create a new migration script in the `db/migrations` folder, with
the format:

	$UNIXTIMESTAMP_$NAME.cql

For example, if you ran `cqlmm create MyNewMigration cql`, the new migration
file will be something like: `db/migrations/1414448037_MyNewMigration.cql`.

The `FORMAT` argument to `cqlmm create` currently supports "cql", only. This
may be extended in the future to support "go", for writing migrations in Go.

Within the new migration, there are two sections, denoted by `-- +cql up`, and
`-- +cql down`. The `-- +cql up` section is where you should write the CQL
statements that will be executed when you run `cqlmm up`. Similarly, you would
put statements in the `-- +cql down` section, to have them executed when you
run `cqlmm down`.

### In practice

cqlmm does not re-order the statements within a migration script, nor does it
inspect the statements; it just executes the statements you put into the file.
So, it is strongly suggested that you write your downgrade statements in the 
opposite order of your upgrade statements. For example:

	-- +cql up
	CREATE TABLE users (
		id uuid,
		username text,
		password text,
		enabled bool,
		PRIMARY KEY (id, username)
	);

	CREATE INDEX users_username_idx ON users (username);

	-- +cql down
	DROP INDEX users_username_idx;
	DROP TABLE users;

### Comments

CQL supports three different comment notations:

*	`--` at the beginning of a line
*	`//` at the beginning of a line
*	`/* ... */` is for multi-line comments, where the opening `/*` is at 
	the beginning of a line

cqlmm only supports the `--` notation, for now.

When cqlmm parses the `.cql` files in the `db/migrations` directory, it skips
any comment lines (lines that begin with `--`), and does not send them to the
server.

For more information about comments in the CQL syntax, please refer to the 
[documentation](https://cassandra.apache.org/doc/cql3/CQL.html#Comments).

Upgrading the schema
--------------------

To apply a new migration to your Cassandra cluster (or node), all you have to
do is run

	$ cqlmm up

This will execute all statements, in all migrations files, that have not yet
been run on the database.

For more information on how cqlmm tracks migrations, there is a page in the
wiki explaining it: 
[How cqlmm tracks migrations](https://github.com/nesv/cqlmm/wiki/How%20cqlmm%20tracks%20migrations)

Downgrading the schema
----------------------

To revert your schema back to a previous version, run

	$ cqlmm down

Unlike `cqlmm up`, `cqlmm down` only reverts by one migration at a time,
although this may be changed, in the future.
