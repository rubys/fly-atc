class AtcGenerator < Rails::Generators::Base
  def generate_app
    source_paths.push File.expand_path("./templates", __dir__)

    ### config/routes.rb

    @routes = IO.read("config/routes.rb")

    unless @routes.include? "fly_atc_scope"
      _, prolog, routes = @routes.split(/(.*Rails.application.routes.draw do\n)/m,2)
      routes, epilog, _ = routes.split(/^(end.*)/m,2)
      routes = routes.split(/\n\s*\n/)
      scoped = routes.select {|route| route =~ /^\s*\w/ && !route.include?('as:')}

      @routes = <<~EOF
        #{prolog.rstrip}
          fly_atc_scope = ENV.fetch("FLY_ATC_SCOPE", "")

          unless fly_atc_scope == ""
            mount ActionCable.server => "/\#{fly_atc_scope}/cable"
          end

          scope fly_atc_scope do
        #{scoped.join("\n\n").gsub(/^ /, "   ")}
          end

        #{(routes-scoped).join("\n\n").rstrip}
        #{epilog.rstrip}
      EOF
    end

    template "routes.erb", "config/routes.rb"

    ### app/views/layouts/application.html.erb

    @layout = IO.read("app/views/layouts/application.html.erb")

    unless @layout.include? "action_cable_meta_tag_dynamic"
      @layout[/<meta.*?\n()\r?\n/m, 1] = "    <%= action_cable_meta_tag_dynamic %>\n"
    end

    template "application.html.erb", "app/views/layouts/application.html.erb"

    ### bin/fly-atc

    unless File.exist?("bin/fly-atc")
      system "bundle binstubs fly-atc"
    end

    ### Dockerfile

    gsub_file "Dockerfile", "bin/thrust", "bin/fly-atc"
  end
end
