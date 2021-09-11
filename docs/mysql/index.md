# MySQL Backup and Restore with Elastio

### Backup databases

Streaming full backups of all databases can be accomplished quite easily

`$ mysqldump -A | elastio stream backup mysql_backup_20210807 --vault default`

### Table level backups

For now, table level restores can only be accomplished by dumping each table to a separate backup.  Get a list of the tables for a given database and loop through them as follows.

```
$ TABLES=$(mysql -e "SHOW TABLES;" <database> | xargs)

$ for table in ${TABLES}; do

> echo "Backing up database.${table} to elastio..."

> mysqldump <database> ${table} | elastio stream backup mysql_DBNAME_${table} --vault default --tag <database>:${table}

> done
```

### Restore a table from a backup

This command will restore a table named `addresses` from the database `customers`.  Note that this will work in conjunction with the naming used above.
```
$ elastio rp list | grep "customers_addresses"
$ elastio stream restore --rp <restore_point_id> > addresses.sql
```

### Alternately, restore directly to the database

```
$ elastio stream restore --rp <restore_point_id> | mysql
```


