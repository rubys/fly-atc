# fly-atc

A SaaS toolkit for converting a personal application into a efficient, siloed, multi-tenant application, where each user of your application is assigned a dedicated virtual machine.

** Work in Progress **

## Usage

This is all TBD at this point, but for Rails projects it is likely to go something like this:

* Replace thruster with fly-atc in Gemfile and Dockerfile
* Define your tenants in a config file, probably YAML.

For non-Rails projects, the process is likely going to be similar:

* Follow the instructions for using thruster with your framework, but substitute fly-atc for thruster.
* Define your tenants in a config file, probably YAML but JSON could also be supported.

Fly.io's dockerfile generators will be able to help with this.

For approximately $1 US per month, you can run:
  * [1 performance machine with 2Gb of RAM, 10GB of bandwidth, and 5GB of storage for 15 hours/month](https://fly.io/calculator?m=0_0_0_0_0&f=c&b=iad.10&a=no_none&r=shared_0_1_iad&t=10_100_5&u=0_1_100&g=1_performance_15_1_2048_iad_1024_0).
  * [1 shared machine with 1Gb of RAM, 10GB of bandwidth, and 5GB of storage for 80 hours/month](https://fly.io/calculator?m=0_0_0_0_0&f=c&b=iad.10&a=no_none&r=shared_0_1_iad&t=10_100_5&u=0_1_100&g=1_shared_80_1_1048_iad_1024_0).

Vertical scaling can be achieved by adding more machines.

## Motivation

I've been running my [Showcase](https://github.com/rubys/showcase?tab=readme-ov-file#showcase) software for nearly three years.  Things have changed over time that I now want to take advantage of.  I want take the opportunity to package those changes in the form of a toolkit that others can take advantage of.

From Wikipedia description of [SaaS](https://en.wikipedia.org/wiki/Software_as_a_service):

> SaaS customers have the abstraction of limitless computing resources, while [economy of scale](https://en.wikipedia.org/wiki/Economy_of_scale) drives down the cost. SaaS architectures are typically [multi-tenant](https://en.wikipedia.org/wiki/Multi-tenant); usually they share resources between clients for efficiency, but sometimes they offer a siloed environment for an additional fee.

The focus of this toolkit is efficient, siloed, multi-tenant applications *with no changes to the application*, taking advantage of:

* [Auto-suspend](https://community.fly.io/t/autosuspend-is-here-machine-suspension-is-enabled-everywhere/20942) -  Virtual Machines that pop into existence when needed and disappear when not in use.
* [SQLite ready for production](https://rubyonrails.org/2024/11/7/rails-8-no-paas-required#getting-sqlite-ready-for-production) - raw performance coupled with operational compression of complexity; see [Supercharge the One Person Framework with SQLite: Rails World 2024](https://fractaledmind.github.io/2024/10/16/sqlite-supercharges-rails/).
* [Litestream](https://litestream.io/) -  No-worry backups.  Virtual machines can be literally destroyed and recreated elsewhere and start back up exactly where they left off.
* [Tigris Global Storage](https://fly.io/docs/tigris/) - globally caching, S3-compatible object storage.

That's a lot of moving parts.  I've documented my [current architecture](https://github.com/rubys/showcase/blob/main/ARCHITECTURE.md) and published a [blueprint](https://fly.io/docs/blueprints/shared-nothing/).

The goal of fly-atc is to enable you configure multiple tenants and then not worry about this further, enabling you to focus on your application.

## Approach

For illustrative purposes consider a SaaS Calender application implemented in Ruby on Rails using SQLite3 as the database.  (My showcase application is a bit more involved than a calendar, but those details aren't important).

Key concepts:

* Each user/customer has a primarly location, and is assigned a single machine near that location.  Such machines can be accessed from anywhere, but have lower latency near that location.
* Each user can have multiple calendars.  Each calendar is associated with a single tenant on the user's machine.  Each tenant consists a running instance of the web server application with one ([or more](https://rubyonrails.org/2024/11/7/rails-8-no-paas-required#a-solid-reduction-of-dependencies)) databases.

With that in mind, consider the following URL paths:

* `/bellevue/2025/winter/`
* `/bellevue/2025/summer-medal-ball/`
* `/bellevue/2025/summer-showcase/`
* `/boston/2025/april/`
* `/boston/2025/mini-comp/`
* `/boston/2025/october/`
* `/livermore/2025/the-music-of-prince/`
* `/livermore/2025/james-bond/`
* `/raleigh/2025/disney/`
* `/raleigh/2025/in-house/`

The first segment of the path identifies the user, and therefore the machine.  The next two segments combined identify the tenant on that machine.  This is but a subset of the planned showcases, you can see a [full list](https://smooth.fly.dev/showcase/) or even a [map](https://smooth.fly.dev/showcase/regions/) (click on the arrows under the map to move to different continents).

`fly-atc`'s responsibilities are to:
* Route requests to the correct machine
* Ensure databases are present/restored from backup
* Start/stop tenants as required
* Hand off requests to tenants

Rails 8 introduces [thruster](https://rubyonrails.org/2024/11/7/rails-8-no-paas-required#enter-kamal-2--thruster).  `fly-atc` is a replacement for thruster:
  * thruster requires no configuration, is limited to a single tenant.
  * fly-atc enables multiple tenants, based on your configuration.

## Implementation

Based on:
* [Thruster](https://github.com/basecamp/thruster) ([announcement](https://dev.37signals.com/thruster-released/))
* [tinyrp](https://github.com/pgaijin66/tinyrp) ([docs](https://prabeshthapa.medium.com/learn-reverse-proxy-by-creating-one-yourself-using-go-87be2a29d1e))

Near term plans:

* Remove certificate/https support
* Add launch on request / shutdown on idle
* Add [fly-replay](https://fly.io/docs/networking/dynamic-request-routing/)

On the radar:

* Support for targets other than fly.io.
* Support for platforms other than Rails, likely starting with Node, and focusing on popular ORMs: [Prisma](https://www.prisma.io/), [TypeORM](https://typeorm.io/), and [Sequelize](https://sequelize.org/).
* Dashboard.  One should be able to deploy new users and make other configuration changes using only your cell phone.  I [do this today](https://github.com/rubys/showcase/blob/main/ARCHITECTURE.md#administration) with my showcase application.