# 1. Install Go

```cmd
    curl -OL https://golang.org/dl/go1.16.7.linux-amd64.tar.gz
    sha256sum go1.16.7.linux-amd64.tar.gz
    sudo tar -C /usr/local -xvf go1.16.7.linux-amd64.tar.gz
    sudo nano ~/.profile
    export PATH=$PATH:/usr/local/go/bin
    source ~/.profile
    go version
```

# 2. Install SQLC

```cmd
    sudo snap install sqlc
```
