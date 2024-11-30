# fly-atc

A SaaS toolkit for converting a personal application into a efficient, siloed, multi-tenant application, where each user of your application is assigned a dedicated virtual machine.

** Work in Progress **

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

 * [Demo](./docs/demo.md)
 * [Overview](./docs/overview.md)
 * [Iaas vs PaaS vs SaaS](./docs/paas.md)
 * [SQlite3](./docs/sqlite3.md)
 * [Todos](./docs/todos.md)
