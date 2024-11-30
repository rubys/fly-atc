Hi, I'm Sam Ruby:

  * Author of [Agile Web Development with Rails 8](https://pragprog.com/titles/rails8/agile-web-development-with-rails-8/).
  * I develop and run [Showcase software](https://github.com/rubys/showcase?tab=readme-ov-file#showcase) for Dance studios at [smooth.fly.dev](https://smooth.fly.dev/).
  * I am an employee of [Fly.io](https://fly.io/about/).

My showcase software has been used in production for almost 3 years.  Originally hosted on my mac-mini in my attic and used for local events in the mid-atlantic states, it now runs on fly.io and is used by 50 dance studios in 5 countries on 3 continents.

I've written up and given presentations on the current architecture: [Shared Nothing](https://fly.io/docs/blueprints/shared-nothing/), and in the spirit of Rails, am in the process of extracting the code into a toolkit that can be used by others.

fly-atc is the result.

It currently is Rails and Fly.io specific.  I'm seeking:

* Users and Feedback on usage with Rails and Fly.io.
* Collaborators to extend this code base to other frameworks and, yes, even other platforms.

## Motivation

Things have changed over time that I now want to take advantage of.  I also want take the opportunity to package those changes in the form of a toolkit that others can take advantage of.

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

`fly-atc`'s responsibilities are to:
* Route requests to the correct machine
* Ensure databases are present/restored from backup
* Start/stop tenants as required
* Hand off requests to tenants

Rails 8 introduces [thruster](https://rubyonrails.org/2024/11/7/rails-8-no-paas-required#enter-kamal-2--thruster).  `fly-atc` is a replacement for thruster:
  * thruster requires no configuration, is limited to a single tenant.
  * fly-atc enables multiple tenants, based on your configuration.

On the radar:

* Support for targets other than fly.io.
* Support for platforms other than Rails, likely starting with Node, and focusing on popular ORMs: [Prisma](https://www.prisma.io/), [TypeORM](https://typeorm.io/), and [Sequelize](https://sequelize.org/).
* Dashboard.  One should be able to deploy new users and make other configuration changes using only your cell phone.  I [do this today](https://github.com/rubys/showcase/blob/main/ARCHITECTURE.md#administration) with my showcase application.