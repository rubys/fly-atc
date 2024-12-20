# Todos and discussion items

Short range projects (thought to be relatively easy, but not necessary for the first proof of concept):

* Applications are currently launched on demand; there should be an option to launch immediately.
* Applications should (possibly optionally) stop when idle.
* Modifying the Rails [flyctl scanner](https://github.com/superfly/flyctl/blob/master/scanner/rails.go) and/or [dockerfile generator](https://github.com/fly-apps/dockerfile-rails?tab=readme-ov-file#overview) to automatically configure litestream for sqlite3 based deployments.
* Generator option to create a [stimulus controller](https://github.com/rubys/showcase/blob/main/app/javascript/controllers/region_controller.js) to insert 
that <a href="https://github.com/rubys/showcase/blob/main/app/javascript/controllers/region_controller.js">listens</a> for
<a href="https://turbo.hotwired.dev/reference/events#http-requests" rel="nofollow">turbo:before-fetch-requests</a> events and
inserts a <a href="https://fly.io/docs/networking/dynamic-request-routing/#the-fly-prefer-region-request-header" rel="nofollow">fly-prefer-region header</a>.
The instance/region is extracted from a data attribute added to the <a href="https://github.com/rubys/showcase/blob/aae08a6d57f92335b2cdbb94756e5416b7b50f83/app/views/layouts/application.html.erb#L16">body</a> tag.
* Simplify configurations: if you specify a name, infer an endpoint and vice versa, from there infer database.
* If there are no routes defined anchored on root, respond to health checks and provide something in response to requests to the root url itself: perhaps a redirect to one of the routes anchored on this machine, chosen at random?

May need more investigation:
* If one ever switches an application back to a machine where it once was running there may be a stale database there.  Restore from litestream when such occurs.
* Would it be better to run without a volume?
* Provide the ability to route based on machine metadata rather than instance id.
* Automatically configure suspend and/or volume expansion?
* What should `fly console` do?  It is relatively straightforward to retrieve but not replicate databases, effectively providing a snapshot of a current database to play with disconnected from production.

Longer range:

* Automatic provisioning - atc.yml is a declarative description of the intended provisioning of tenants; make use of the machine API to start machines on demand.
* Automatic failover; if a machine goes down, start a new one.
* Other frameworks and ORMs; probably starting with [Prisma](https://www.prisma.io/)
* Instead of a config file in the source tree, perhaps the configuration should be in some sort of key/value store allowing dynamic updates
* Some sort of dashboard allowing the configuration to be updated using a web interface
