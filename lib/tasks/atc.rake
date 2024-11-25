if ENV["FLY_ATC_SCOPE"]
  actions = Rake::Task["db:prepare"].actions.clone

  namespace :db do
    task :atc_prepare => "db:load_config" do
      actions.each {|action| action.call}
    end
  end

  Rake::Task["db:prepare"].clear
end
