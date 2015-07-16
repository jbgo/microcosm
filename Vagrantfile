# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "private_network", ip: "192.168.33.33"
  config.vm.synced_folder "./", "/home/vagrant/go/src/github.com/jbgo/microcosm", owner: "vagrant"
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
  end

  config.vm.provision "shell", path: "provision.sh"
end
