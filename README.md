# sqltools
An under-development suite of tools for working with different types of SQL database types, including, but not necessarily limited to, MySQL, MariaDB, Postgres, Microsoft™ SQL Server®, maybe a little Oracle with maybe some SQLite thrown in for good measure.

## What's Included?
- `mssqldump` - My attempt at mimicing the `mysqldump` command against MS SQL Server
- `mysqltablerestore` - Restore a single table from a mysql dump

## Known Limitations or Defects
- `mssqldump` provides support for the most commonly-used objects and attributes.  I'm guessing these would cover about 80% of all cases, though I really just made that number up right now
- `mysqltablerestore` - It is incredibly dumb, works only for tables and assumes usage of `mysqldump` to create the dump in the first place

## TODO List
- [ ] Better name for the package
- [ ] Tests