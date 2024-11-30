### Preface: Clarke's three laws

> 1. When a distinguished but elderly scientist states that something is possible, he is almost certainly right. When he states that something is impossible, he is very probably wrong.
> 2. The only way of discovering the limits of the possible is to venture a little way past them into the impossible.
> 3. Any sufficiently advanced technology is indistinguishable from magic.

### Prerequisites

* Ruby >= 3.2.0.  Available via [brew](https://formulae.brew.sh/formula/ruby) and [mise](https://mise.jdx.dev/lang/ruby.html).
* [flyctl](https://fly.io/docs/flyctl/install/)

Windows users are strongly encouraged to use [WSL](https://learn.microsoft.com/en-us/windows/wsl/about).

---

# ATC Demo

In an empty directory run:

```bash
fly launch --region iad --from https://github.com/rubys/depot8.git
fly apps open
```

At this point you have a SaaS: a product catalog, as defined by chapter 6 in [Agile Web Development with Rails 8](https://pragprog.com/titles/rails8/agile-web-development-with-rails-8/).

If you have two browser windows open, you can see a list of products in one, and make a change in the other and see the product list update live.

Key features:

* Active Record (Sqlite3)
* Active Storage (S3/Tigris)
* Active Job (Solid Queue)
* Action Cable (Web Sockets)

Now lets add Litestream:

```bash
bundle add fly-atc
bin/rails generate atc
fly deploy
```

Accept the changes proposed by the generator.

At this point you can verify that Litestream is active by looking at the logs or looking into your Tigris bucket in the fly.io dashboard.

```bash
fly logs
fly dashboard
```

We've now got a single database with a single tenant on a single machine.  Now lets go to multiple tenants.  Create a `config/atc.yaml` with the following contents:

```yaml
routes:
  - name: Sam
    endpoint: /sam
    database: sam
  - name: Ben
    endpoint: /ben
    database: ben
  - name: Annie
    endpoint: /annie
    database: annie
  - name: Darla
    endpoint: /darla
    database: darla
  - name: index
    database: production
```

Note the "extra" tenant at the bottom.  We will use it to display an index.  Modify `config/routes.rb` to make this happen:


```diff
  scope fly_atc_scope do
    get "atc/index"
    resources :products
    # Define your application routes per the DSL in 
    # https://guides.rubyonrails.org/routing.html

    # Defines the root path route ("/")
+   if fly_atc_scope == ""
+     root "atc#index"
+   else
      root "products#index"
+   end
  end
```

Now run:

```bash
fly deploy
```

Visit any store other than Sam's.  Make any change you like (perhaps change a price)

Now lets add regions to all but the index (which will run in all regions):

```yaml
routes:
  - name: Sam
    endpoint: /sam
    database: sam
    region: iad
  - name: Ben
    endpoint: /ben
    database: ben
    region: den
  - name: Annie
    endpoint: /annie
    database: annie
    region: sea
  - name: Darla
    endpoint: /darla
    database: darla
    region: fra
  - name: index
    database: production
```

Let's create three more machines:

```bash
fly scale count 3 --region=den,sea,fra
fly deploy
```

Visit the store where you made the change.  Verify that the change is present.  Visit another store and see that the change is not there.  Check the logs to verify that each store is being served by a machine in the correct region.

Close the all browser windows and wait a few minutes, and watch all the machines automatically stop.

Optional things to explore:

* When you scale to the point where you want more than one machine in a region, you can specify a machine `instance` instead of a region in `atc.yaml`.
* Change _auto_stop_machines = 'stop'` to `auto_stop_machines = 'suspend'` in `fly.toml`.  Redeploy.  After this change, machines will wake up faster.
* In the `[[mounts]]` section add the following lines (adjusting to taste) to make the volumes automatically grow:

```toml
  auto_extend_size_threshold = 80
  auto_extend_size_increment = "1GB"
  auto_extend_size_limit = "100GB"
```