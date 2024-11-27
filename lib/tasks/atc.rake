actions = Rake::Task["db:prepare"].actions.clone

namespace :db do
  task :atc_prepare => "db:load_config" do
    actions.each {|action| action.call}
  end
end

Rake::Task["db:prepare"].clear

namespace :litestream do
  task :atc_config => "db:load_config" do
    require 'erubi'

    @dbs =
      ActiveRecord::Base
        .configurations
        .configs_for(env_name: "production", include_hidden: true)
        .select { |config| ["sqlite3", "litedb"].include? config.adapter }
        .map(&:database)

    @config  = ENV["LITESTREAM_CONFIG"] || Rails.root.join("config/litestream.yml")

    template = File.read(File.join(File.dirname(__FILE__), "templates/litestream.yml.erb"))
    result = eval(Erubi::Engine.new(template).src)

    unless File.exist?(@config) && File.read(@config) == result
      File.write(@config, result)
    end
  end

  task :atc_restore => "litestream:atc_config" do
    next unless ENV["BUCKET_NAME"]

    @dbs.each do |database|
      next if File.exist? database
      sh "bundle exec litestream restore -config #{@config} -if-replica-exists #{database}"
    end
  end

  task :atc_replicate => "litestream:atc_config" do
    next unless ENV["BUCKET_NAME"]
    sh "bundle exec litestream replicate -config #{@config}"
  end
end

namespace :atc do
  task :prepare => ["litestream:atc_restore", "db:atc_prepare"]
end
