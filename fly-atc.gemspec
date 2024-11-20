require_relative "lib/fly-atc/version"

Gem::Specification.new do |s|
  s.name        = "fly-atc"
  s.version     = FlyAtc::VERSION
  s.summary     = "A SaaS toolkit"
  s.description = "An HTTP/2 proxy for mutli-tenant production deployments"
  s.authors     = [ "Sam Ruby" ]
  s.email       = "rubys@intertwingly.net"
  s.homepage    = "https://github.com/rubys/fly-atc"
  s.license     = "MIT"

  s.metadata = {
    "homepage_uri" => s.homepage
  }

  s.files = Dir[ "{lib}/**/*", "MIT-LICENSE", "README.md" ]
  s.bindir = "exe"
  s.executables << "fly-atc"
end
