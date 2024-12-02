# fly-atc

* A SaaS toolkit for converting a personal application into a efficient, siloed, multi-tenant application, where each user of your application is assigned a dedicated virtual machine.
* Zero-config no-worry sqlite3 backups using Litestream

Currently supports Rails and fly.io; looking for collaborators to expand to other frameworks and platforms.

*** **Work in Progress** ***

Do not use in production just yet.

## Usage

### Quickstart (single tenant):

```
bundle add fly-atc
bundle binstubs fly-atc
```

Replace `thruster` with `fly-atc` in Dockerfile.

### Quickstart (multi-tenant):

```
bundle add fly-atc
bin/rails generate atc
```

Edit `config/atc.yml` as needed.

Vertical scaling can be achieved by adding more machines.

## More information:

 * [Overview](./docs/overview.md) - why this toolkit was created.
 * [Demo](./docs/demo.md) - up and running in minutes.
 * [Config](./docs/config.md) - configuration options
 * [Iaas vs PaaS vs SaaS](./docs/paas.md) - Rails never needed a PaaS; Rails needs a PaaS now more than ever.
 * [SQlite3](./docs/sqlite3.md) - perhaps Sqlite3 isn't right for you.
 * [Todos](./docs/todos.md) - where we go from here.
