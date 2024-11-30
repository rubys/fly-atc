require "bundler/setup"
require "bundler/gem_tasks"

task :release_all => :release do
  Dir.chdir("pkg") do
    Dir.glob("*-#{FlyAtc::VERSION}-*.gem").each do |gem|
      sh "gem push #{gem}"
    end
  end
end
