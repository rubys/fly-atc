module FlyAtc
end

require_relative "fly-atc/version"
require_relative "helpers/atc-cable"

class FlyAtcRailtie < Rails::Railtie
  rake_tasks do
    Dir[File.join(File.dirname(__FILE__),'tasks/*.rake')].each { |f| load f }
  end
end
