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
	
