# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure('2') do |config|
  # base box
  config.vm.box = 'bento/ubuntu-18.04'

  # forward ports
  config.vm.network 'forwarded_port', guest: 5672, host: 5672 # rabbit port
  config.vm.network 'forwarded_port', guest: 15672, host: 15672 # rabbit management port

  # customize the VM
  config.vm.provider 'virtualbox' do |v|
    v.cpus = 1 # use one processor
    v.memory = 1024 # use 1024 of RAM
  end

  # provision the VM
  config.vm.provision 'shell', path: 'vagrant.sh', privileged: false, keep_color: true
end
