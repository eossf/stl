# STL : SainTeLyon flutter app
## Installation from ZERO - Debian 10
````bash
apt update
apt -y install git
git clone https://github.com/eossf/common.git
cd common
find . -name "*.sh" -exec chmod 775  {} \;
````
### Create stl namespace and registry on port 5000
````bash
cd scripts
./install_debian10.sh stl 5000

# change namespace
kubectl config set-context --current --namespace=stl
````
### Useful aliases for kubectl
````bash
cat <<EOF >> ~/.bashrc
export LS_OPTIONS='--color=auto'
eval "`dircolors`"
alias ls='ls $LS_OPTIONS'
alias ll='ls $LS_OPTIONS -l'
alias l='ls $LS_OPTIONS -ltrhA'
source <(kubectl completion bash)
complete -F __start_kubectl k
alias k='kubectl '
alias kpo='kubectl get pods'
alias klo='kubectl logs -f'
alias kcm='kubectl get cm'
alias ksv='kubectl get svc'
alias kde='kubectl get deployments'
alias kns='kubectl get ns'
EOF
source ~/.bashrc
````
### Install Helm mongodb from bitnami to the current context/ns
````bash
helm repo add bitnami https://charts.bitnami.com/bitnami

# node port
helm install --set image.tag=3.6.23 --set service.type="NodePort" --set service.nodePort=30000 db-stl bitnami/mongodb

# ClusterIp
helm install --set image.tag=3.6.23 db-stl bitnami/mongodb

````

````bash
# (...)
#MongoDB(R) can be accessed on the following DNS name(s) and ports from within your cluster:
export MONGODB_CLUSTER_DNS="db-stl-mongodb.default.svc.cluster.local"

# get IP of the standalone server (CIDR > 16), first interface !
vlan16=`ip r | grep "default" | cut -d" " -f3 | cut -d"." -f1-2`
export MONGO_HOST=`ip -o -4 addr list | grep "$vlan16" | awk '{print $4}' | cut -d/ -f1 | head -1`

#To get the root password run:
export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace stl db-stl-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

#To connect to your database, create a MongoDB(R) client container:
kubectl run --namespace stl db-stl-mongodb-client --rm --tty -i --restart='Never' --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:3.6.23 --command -- bash

#Then, run the following command:
mongo admin --host "db-stl-mongodb" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD

#To connect to your database from outside the cluster execute the following commands:
# direction connection ClusterIp
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb 27017:27017 &
mongo --host 127.0.0.1 --port 27017 --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD

# or, with NodePort
export NODE_IP=$(kubectl get nodes --namespace stl -o jsonpath="{.items[0].status.addresses[0].address}")
export NODE_PORT=$(kubectl get --namespace stl -o jsonpath="{.spec.ports[0].nodePort}" services mongodb-stl)
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb 27017:$NODE_PORT &
mongo --host $NODE_IP --port $NODE_PORT --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD
````
### Install noSqlclient
````bash
docker run -d -p 3000:3000 --name mongoclient -e MONGOCLIENT_DEFAULT_CONNECTION_URL="mongodb://root:$MONGODB_ROOT_PASSWORD@$MONGO_HOST/admin?ssl=false" -e MONGOCLIENT_AUTH="true" -e MONGOCLIENT_USERNAME="root" -e MONGOCLIENT_PASSWORD="$MONGODB_ROOT_PASSWORD" mongoclient/mongoclient:latest
ufw allow 27017
````

### Feed with data
````bash
# java script
db = db.getSiblingDB("stl");
db.createUser({
  user: "stluser",
  pwd: "<PASSWORD>",
  roles: [{role: "readWrite", db: "stl"}, { role: "dbAdmin", db: "stl" } ]
});
db.createCollection("track");
db.track.insertOne({"id": "1", "name": "Track Test ID 1"});

# MongoShell command
mongo mongodb://stluser:<PASSWORD>@db-stl-mongodb:27017/stl

# useful commands
db.track.find();
    { "_id" : ObjectId("60a62194e75136aca562d2bc"), "ID" : 0, "Name" : "Track Test ID 1" }
````

# TODO
 - [x] how to stop kubectl port-forward --namespace stl svc/db-stl-mongodb 27017:27017
 - [x] launch docker mongoclient with --name
 - [ ] traefik 27017
 - [x] create db stl + collection + user
 - [ ] index + minimal data ==> initConfigMap.name 	Custom config map with init scripts 	nil

# Revert Back
````bash
# port forwarding if apply
kill -9 `ps -aux | grep "27017:27017" | grep "kubectl" | awk '{print $2}'`
# revert back mongo
db.dropUser("stluser");
db.track.drop();
db.dropDatabase();
````


