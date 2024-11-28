# Todos and discussion items

Short range projects (thought to be relatively easy, but not necessary for the first proof of concept):

* Applications are currently launched on demand; there should be an option to launch immediately.
* Applications should (possibly optionally) stop when idle.
* Router handler should reverse proxy if the payload exceeds 1Mb (essentially any non-GET, non-HEAD request that either doesn't have a Content-length header, or has one that specified a content that is too large).  Reverse proxying to the same URL but with an added `fly-force-instance-id` or `fly-prefer-region` header should do the trick.
* Modifying the Rails [flyctl scanner](https://github.com/superfly/flyctl/blob/master/scanner/rails.go) and/or [dockerfile generator](https://github.com/fly-apps/dockerfile-rails?tab=readme-ov-file#overview) to automatically configure litestream for sqlite3 based deployments.
* Generator option to create a [stimulus controller](https://github.com/rubys/showcase/blob/main/app/javascript/controllers/region_controller.js) to insert 
that <a href="https://github.com/rubys/showcase/blob/main/app/javascript/controllers/region_controller.js">listens</a> for
<a href="https://turbo.hotwired.dev/reference/events#http-requests" rel="nofollow">turbo:before-fetch-requests</a> events and
inserts a <a href="https://fly.io/docs/networking/dynamic-request-routing/#the-fly-prefer-region-request-header" rel="nofollow">fly-prefer-region header</a>.
The instance/region is extracted from a data attribute added to the <a href="https://github.com/rubys/showcase/blob/aae08a6d57f92335b2cdbb94756e5416b7b50f83/app/views/layouts/application.html.erb#L16">body</a> tag.

May need more investigation:
* If one ever switches an application back to a machine where it once was running there may be a stale database there.  Restore from litestream when such occurs.
* Would it be better to run without a volume?

Longer range:

* Automatic provisioning
* Automatic failover
* Other frameworks and ORMs; probably starting with [Prisma](https://www.prisma.io/)