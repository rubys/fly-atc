Be sure that sqlite3 is right for you - it might not be!

### Background information

[Supercharge the One Person Framework with SQLite: Rails World 2024](https://fractaledmind.github.io/2024/10/16/sqlite-supercharges-rails/)

* SQLite3 "out of the box" is not suited for production, but by fine tuning the configuration and usage it is.  [Rails 8](https://rubyonrails.org/2024/11/7/rails-8-no-paas-required) does both.
* When running SQLite3 in production, you need to have a solid backup mechanism setup.
* Sqlite3 only supports linear writes, but performs 10 to 600 times faster than Postgres.

### Recommendations

* For single machine, single tenant, applications sqlite3 is an excellent choice in that it provides excellent vertical scalability.  For backups, use [litestream](https://litestream.io/); if/when a machine or volume fails simply start up a new one.  This provides failover, but not automatic failover.
* For multi-tenant machines with readily partitionable data stores, horizontal scalability can be achieved by dynamically routing requests.  This can also improve responsiveness by placing applications near users.
* For applications with data that is not readily partitionable, where reads dominate write requests, consider
[LiteFS](https://fly.io/docs/litefs/).  A [built in HTTP proxy](https://fly.io/docs/litefs/proxy/) is provided for web applications; background jobs require [write forwarding](https://github.com/superfly/litefs/issues/56).
* For all other uses, [PostgreSQL](https://www.postgresql.org/) is recommended.
