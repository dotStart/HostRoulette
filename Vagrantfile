# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.synced_folder ".", "/usr/src"
  config.vm.network "forwarded_port", guest: 6379, host: 6379, host_ip: "127.0.0.1"

  config.vm.provider "virtualbox" do |v|
    v.memory = 4096
    v.cpus = 2
  end

  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y redis
    sed -i.bak -e "s/^bind .*$/bind 0.0.0.0/gm" /etc/redis/redis.conf
    systemctl enable redis-server
    systemctl start redis-server
  SHELL
end
