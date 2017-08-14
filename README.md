# scheduler

## prepare enviorenment

### install go env
1.install golang 1.8 on centos7.3

wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz

tar -xzf go1.8.3.linux-amd64.tar.gz

mv go /usr/local

export GOROOT=/usr/local/go

export GOPATH=$HOME/Projects/Proj1

export GOPATH=$HOME/Projects/Proj1

### install go packages
go get gopkg.in/olivere/elastic.v5

go get github.com/lib/pq

go get github.com/satori/go.uuid

### install liteide
1. download

https://sourceforge.net/projects/liteide

2. unzip tar

tar jxvf litexxxxx.tar.bz2

3. install qt env

yum install qt-x11

4. move liteide to /usr/local

mv liteide /usr/local

5. run liteide

cd /usr/local/liteide/bin

./liteide PATH_TO_YOUR_GO_PROJECT

### install elasticsearch 5.5.1
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.5.1.rpm

sudo rpm --install elasticsearch-5.5.1.rpm

### install postgresql 9.6
https://www.postgresql.org/download/linux/redhat/

configure

http://www.jianshu.com/p/7e95fd0bc91a

## test scheduler
### install python3 on centos7
pip3 install --upgrade pip

pip3 install psycopg2

pip3 install requests