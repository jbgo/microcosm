#!/bin/bash

set -e

echo "===== Installing golang ====="
if [ -f /usr/local/go/bin/go ]; then
  echo " ---> go already installed"
else
  download_path=/usr/local/go-1.4.2.tar.gz

  if ! [ -f $download_path ]; then
    echo ' ---> downloading go-1.4.2'
    wget -qO- https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz > $download_path
  else
    echo ' ---> go-1.4.2 already downloaded'
  fi
  echo ' ---> extracting archive to /usr/local/go'
  tar -C /usr/local -xzf $download_path

  echo 'export GOROOT=/home/vagrant/go' >> /home/vagrant/.profile
  echo " ---> export GOROOT=$GOROOT"

  echo 'export PATH=$PATH:/usr/local/go/bin:$GOROOT/bin' >> /home/vagrant/.profile
  echo " ---> export PATH=$PATH"
fi

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

cert_path=/vagrant/.cert
if [ -f $cert_path/cert.pem ]; then
  echo " ---> TLS keys already generated"
else
  ca_passphrase=notsosecret
  ca_domain=docker
  alt_names="IP:192.168.33.100,IP:127.0.0.1"

  echo " ---> generating docker CA"
  openssl genrsa -aes256 -passout "pass:$ca_passphrase" -out ca-key.pem 2048 2>/dev/null
  openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem \
    -passin "pass:$ca_passphrase" -subj "/CN=$ca_domain"

  echo " ---> generating server keypair"
  openssl genrsa -out server-key.pem 2048 2>/dev/null
  openssl req  -sha256 -new -key server-key.pem -out server.csr -subj \
    "/CN=$ca_domain"
  echo "subjectAltName = $alt_names" > extfile-server.cnf
  openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem \
    -extfile extfile-server.cnf -CAcreateserial -out server-cert.pem \
    -passin "pass:$ca_passphrase"

  echo " ---> generating client keypair"
  openssl genrsa -out key.pem 2048 2>/dev/null
  openssl req -subj '/CN=client' -new -key key.pem -out client.csr
  echo "extendedKeyUsage = clientAuth" > extfile-client.cnf
  openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem \
    -CAcreateserial -out cert.pem -extfile extfile-client.cnf \
    -passin "pass:$ca_passphrase"

  echo " ---> setting permissions on keys"
  mkdir -p $cert_path
  mv *.pem $cert_path/
  chmod 0400 $cert_path/*key.pem
  chmod 0444 $cert_path/ca.pem $cert_path/*cert.pem
fi

echo " ---> configure and restart docker daemon"
# TODO
