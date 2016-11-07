require 'fileutils'
require 'os'

MRUBY_VERSION="1.2.0"

file :mruby do
  #sh "git clone --depth=1 https://github.com/mruby/mruby"
  sh "curl -L -L --fail --retry 3 --retry-delay 1 https://github.com/mruby/mruby/archive/1.2.0.tar.gz -s -o - | tar zxf -"
  FileUtils.mv("mruby-1.2.0", "mruby")
end

APP_NAME=ENV["APP_NAME"] || "fd"
APP_ROOT=ENV["APP_ROOT"] || Dir.pwd
bin_path="#{APP_ROOT}/bin"
tools_path="#{ENV["TOOLS_PATH"]}"

desc "compile binary"
task :compile do
  sh "rm -r tools" unless !Dir.exists?(tools_path)
  sh "mkdir tools" unless Dir.exists?(tools_path)
  sh "mkdir bin" unless Dir.exists?(bin_path)
=begin
  Dir.chdir("bin")
  sh "rm -r win64"
  sh "rm -r win32"
  sh "rm -r mac64"
  sh "rm -r mac32"
  sh "rm -r linux64"
  sh "rm -r linux32"
  sh "mkdir win64" unless Dir.exists?("/go/bin/win64")
  sh "mkdir win32" unless Dir.exists?("/go/bin/win32")
  sh "mkdir mac64" unless Dir.exists?("/go/bin/mac64")
  sh "mkdir mac32" unless Dir.exists?("/go/bin/mac32")
  sh "mkdir linux64" unless Dir.exists?("/go/bin/linux64")
  sh "mkdir linux32" unless Dir.exists?("/go/bin/linux32")
  Dir.chdir("/go/bin/win64")
  Dir.chdir("/go/bin/win32")
  Dir.chdir("/go/bin/mac64")
  Dir.chdir("/go/bin/mac32")
  Dir.chdir("/go/bin")
  Dir.chdir("/go/bin/linux32")
=end
  Dir.chdir("/go/bin")
  if OS.linux?
    if OS.bits == 64
      sh"GOOS=linux GOARCH=amd64 go build /go/source/ff.go"
      Dir.chdir(tools_path)
      sh "curl -L https://github.com/appPlant/ff/releases/download/v#{ENV["FF_VER"]}/ff-#{ENV["FF_VER"]}-x86_64-pc-linux-gnu.tgz  | tar xz"
    elsif OS.bits == 32
      sh"GOOS=linux GOARCH=386 go build /go/source/ff.go"
      Dir.chdir(tools_path)
      sh "curl -L https://github.com/appPlant/ff/releases/download/v#{ENV["FF_VER"]}/ff-#{ENV["FF_VER"]}-i686-pc-linux-gnu.tgz  | tar xz"
    end
  elsif OS.mac?
    if OS.bits == 64
      sh"GOOS=darwin GOARCH=amd64 go build /go/source/ff.go"
      Dir.chdir(tools_path)
      sh "curl -L https://github.com/appPlant/ff/releases/download/v#{ENV["FF_VER"]}/ff-#{ENV["FF_VER"]}-x86_64-apple-darwin14.tgz  | tar xz"
    elsif OS.bits == 32
      sh"GOOS=darwin GOARCH=386 go build /go/source/ff.go"
      Dir.chdir(tools_path)
      sh "curl -L https://github.com/appPlant/ff/releases/download/v#{ENV["FF_VER"]}/ff-#{ENV["FF_VER"]}-i386-apple-darwin14.tgz | tar xz"
    end
  elsif OS.windows?
    if OS.bits == 64
      sh"GOOS=windows GOARCH=amd64 go build /go/source/ff.go"
      Dir.chdir(tools_path)
      sh "curl -L https://github.com/appPlant/ff/releases/download/v#{ENV["FF_VER"]}/ff-#{ENV["FF_VER"]}-x86_64-w64-mingw32.tgz  | tar xz"
    elsif OS.bits == 32
      sh"GOOS=windows GOARCH=386 go build /go/source/ff.go"
      Dir.chdir(tools_path)
      sh "curl -L https://github.com/appPlant/ff/releases/download/v#{ENV["FF_VER"]}/ff-#{ENV["FF_VER"]}-i686-w64-mingw32.tgz  | tar xz"
    end
  end
  sh "echo #{OS.windows?}"
  sh "echo #{OS.mac?}"
  sh "echo #{OS.linux?}"
  sh "echo #{OS.bits}"
end

namespace :test do
  desc "run mruby & unit tests"
  # only build mtest for host
  task :mtest => :compile do
    # in order to get mruby/test/t/synatx.rb __FILE__ to pass,
    # we need to make sure the tests are built relative from mruby_root
    MRuby.each_target do |target|
      # only run unit tests here
      target.enable_bintest = false
      run_test if target.test_enabled?
    end
  end

  def clean_env(envs)
    old_env = {}
    envs.each do |key|
      old_env[key] = ENV[key]
      ENV[key] = nil
    end
    yield
    envs.each do |key|
      ENV[key] = old_env[key]
    end
  end

  desc "run integration tests"
  task :bintest => :compile do
    MRuby.each_target do |target|
      clean_env(%w(MRUBY_ROOT MRUBY_CONFIG)) do
        run_bintest if target.bintest_enabled?
      end
    end
  end
end

desc "run all tests"
Rake::Task['test'].clear
task :test => ["test:mtest", "test:bintest"]

desc "cleanup"
task :clean do
  sh "rake deep_clean"
end
=begin
desc "generate a release tarball"
task :release => :compile do
  require 'tmpdir'

  Dir.chdir(mruby_root) do
    # since we're in the mruby/
    release_dir  = "releases/v#{APP_VERSION}"
    release_path = Dir.pwd + "/../#{release_dir}"
    app_name     = "#{APP_NAME}-#{APP_VERSION}"
    FileUtils.mkdir_p(release_path)

    Dir.mktmpdir do |tmp_dir|
      Dir.chdir(tmp_dir) do
        MRuby.each_target do |target|
          next if name == "host"

          arch = name
          bin  = "#{build_dir}/bin/#{exefile(APP_NAME)}"
          FileUtils.mkdir_p(name)
          FileUtils.cp(bin, name)

          Dir.chdir(arch) do
            arch_release = "#{app_name}-#{arch}"
            puts "Writing #{release_dir}/#{arch_release}.tgz  | tar xz"
            `tar czf #{release_path}/#{arch_release}.tgz  | tar xz *`
          end
        end

        puts "Writing #{release_dir}/#{app_name}.tgz  | tar xz"
        `tar czf #{release_path}/#{app_name}.tgz  | tar xz *`
      end
    end
  end
end

namespace :local do
  desc "show version"
  task :version do
    puts APP_VERSION
  end
end

def is_in_a_docker_container?
  `grep -q docker /proc/self/cgroup`
  $?.success?
end

Rake.application.tasks.each do |task|
  next if ENV["MRUBY_CLI_LOCAL"]
  unless task.name.start_with?("local:")
    # Inspired by rake-hooks
    # https://github.com/guillermo/rake-hooks
    old_task = Rake.application.instance_variable_get('@tasks').delete(task.name)
    desc old_task.full_comment
    task old_task.name => old_task.prerequisites do
      abort("Not running in docker, you should type \"docker-compose run <task>\".") \
        unless is_in_a_docker_container?
      old_task.invoke
    end
  end
end
=end

