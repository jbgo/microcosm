#!/bin/bash

set -e

home=/home/vagrant
goroot=/usr/local/go
gopath="$home/go"
project_repo=github.com/jbgo/microcosm
project_src="$gopath/src/$project_repo"
go_version=1.4.2
go_download_url=https://storage.googleapis.com/golang/go$go_version.linux-amd64.tar.gz

if ! [ -x /usr/bin/git ]; then
  echo "==== Installing dependencies ====="
  echo " ---> installing git"
  apt-get update -y && apt-get install -y git-core
fi

echo "===== Installing golang $go_version ====="
if [ -f $goroot/bin/go ]; then
  echo " ---> go already installed"
else
  download_path=$goroot-$go_version.tar.gz

  if ! [ -f $download_path ]; then
    echo " ---> downloading go-$go_version"
    wget -qO- $go_download_url > $download_path
  else
    echo " ---> go-$go_version already downloaded"
  fi
  echo " ---> extracting archive to $goroot"
  tar -C /usr/local -xzf $download_path
fi

echo "===== Configuring golang ====="

profile=/home/vagrant/.profile
if ! grep GOPATH=$gopath $profile; then
  echo " ---> updating $profile"
  echo "export PATH=\$PATH:$goroot/bin:$gopath/bin" >> $profile
  echo "export GOROOT=$goroot" >> $profile
  echo "export GOPATH=$gopath" >> $profile
  echo "cd $project_src" >> $profile
fi
echo " ---> GOPATH=$gopath GOROOT=$goroot"

chown -R vagrant:vagrant $home

echo ' ---> installing godep'
su - vagrant -c 'go get github.com/tools/godep'

echo ' ---> restoring project dependencies'
su - vagrant -c "cd $project_src && pwd && GOPATH=$gopath godep restore"

echo "===== Installing docker ====="
if [ -f /usr/bin/docker ]; then
  echo " ---> docker already installed"
else
  echo " ---> installing latest docker using official installation script"
  wget -qO- https://get.docker.com/ | sh
fi

echo "===== Configuring docker ====="
echo " ---> adding vagrant user to docker group"
usermod -aG docker vagrant

cert_path=$project_src/.cert
if [ -f $cert_path/cert.pem ]; then
  echo " ---> TLS keys already generated"
else
  ca_passphrase=notsosecret
  ca_domain=vagrant-microcosm
  alt_names="IP:192.168.33.33,IP:127.0.0.1"

  echo " ---> generating docker CA"
  openssl genrsa -aes256 -passout "pass:$ca_passphrase" -out ca-key.pem 2048 2>/dev/null
  openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem \
    -passin "pass:$ca_passphrase" -subj "/CN=$ca_domain"

  echo " ---> generating server keypair"
  openssl genrsa -out server-key.pem 2048 2>/dev/null
  openssl req  -sha256 -new -key server-key.pem -out server.csr -subj \
    "/CN=$ca_domain"
  echo "subjectAltName = $alt_names" > extfile.cnf
  openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem \
    -CAcreateserial -extfile extfile.cnf -CAcreateserial -out server-cert.pem \
    -passin "pass:$ca_passphrase"

  echo " ---> generating client keypair"
  openssl genrsa -out key.pem 2048 2>/dev/null
  openssl req -subj '/CN=client' -new -key key.pem -out client.csr
  echo "extendedKeyUsage = clientAuth" > extfile.cnf
  openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem \
    -CAcreateserial -out cert.pem -extfile extfile.cnf \
    -passin "pass:$ca_passphrase"

  echo " ---> setting permissions on keys"
  rm *.csr *.cnf
  mkdir -p $cert_path
  mv *.pem $cert_path/
  chmod 0400 $cert_path/*key.pem
  chmod 0444 $cert_path/ca.pem $cert_path/*cert.pem
fi

echo " ---> configure and restart docker daemon"
tls_opts="--tlscacert=$cert_path/ca.pem --tlskey=$cert_path/server-key.pem --tlscert=$cert_path/server-cert.pem"
echo "DOCKER_OPTS='-D -H unix:// --tlsverify $tls_opts -H=0.0.0.0:2376'" > /etc/default/docker
service docker restart
