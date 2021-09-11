
# Postgres Backup and Restore with Elastio

In its most basic form, the `pg_dump` command creates a database dump as an SQL script to be executed with the `psql` command.  This is perfect for restoration of an entire database, however it may be desireable to pull certain tables out of a backup.

### Perform a full back up of `postgres` database

This will create a backup which we can use to restore in further examples.

```
# Set SCALEZ_STOR_URL to keep from having to specify it each time
<<<<<<< Updated upstream
$ pg_dump postgres | elastio stream backup my_pg_dump_20210807 --vault default --tag postgres:20210807
=======
$ pg_dump postgres | elastio stream backup my_pg_dump_20210807
>>>>>>> Stashed changes
```


### Restore full back up from stream

Use the restore point id for this backup/

```
$ elastio rp list --tag 20210807
$ elastio stream restore --rp <restore_point_id> | psql
```

### Custom-format Backup to Use `pg_restore` Features 

```
$ pg_dump -Fc postgres | elastio stream backup my_pg_dump_20210807 --vault default

# Later restore a single table
# Create a temp database
$ psql -c "CREATE DATABASE `temp`;" postgres

# Then restore it
$ elastio stream restore --rp <restore_point_id> | pg_restore -t users -d temp
```

### Restore Single Table From a Normal Dump

This restores a normal SQL dump, converts it to a custom format dump, then restores one table to a file

```
$ elastio stream restore --rp <restore_point_id> | psql -d tempdb | pg_dump -Fc tempdb | pg_restore -t <my_table> -d <db_name>
```
