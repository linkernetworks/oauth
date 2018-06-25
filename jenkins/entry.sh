#!/bin/bash

UID=`id -u`
GID=`id -g`

echo 'jenkins:x:'$GID':' >> /etc/group
echo 'jenkins:x:'$UID':'$GID':,,,:/home/jenkins:/bin/bash' >> /etc/passwd
sudo chown -R jenkins:jenkins /home/jenkins

##
## Start Docker daemon
##
nohup dockerd 2>&1 > /dev/null &

##
## Start MongoDN daemon
##
mkdir -p /home/jenkins/data/mongo
nohup mongod --dbpath=/home/jenkins/data/mongo 2>&1 > /dev/null &

bash
