module FlyAtc
end

require_relative "fly-atc/version"
require_relative "helpers/atc-cable"

class FlyAtcRailtie < Rails::Railtie
  initializer "fly-atc.configure_rails_initialization" do
    Rails.application.load_tasks

    Rails.configuration.after_initialize do
      Rake::Task['db:atc_prepare'].invoke
    end
  end

  rake_tasks do
    Dir[File.join(File.dirname(__FILE__),'tasks/*.rake')].each { |f| load f }
  end
end
