# STL : SainTeLyon flutter app
STL is a demo of a fullstack application, features are :
- [x] Live Run, follows an existing track (GPX)
- [ ] Live Run, records a track (GPX)
- [x] Support a runner: go to stops using GPS mobile
- [ ] Live View runners and position in a map
- [ ] Guide a runner
- [ ] Misc: Import, Export GPX, Enrich with metadata, ...

Clone this project on a remote server
````sh
ssh -i ~/.ssh/id_rsa root@REMOTE_IP
apt -y install git
git clone https://github.com/eossf/stl.git
````
## Automatic installation
This script read "bash" code blocks in this file. "sh" blocks are ignored. 
````sh
cd stl
scripts/mbe.sh readme.md
````
## Manually installation
````bash
git clone https://github.com/eossf/common.git
cd common
# do this on the git folder: find . -name "*.sh" -exec git add --chmod=+x {} \;
````
### Full k3s/k3d installation
After this step the kubernetes k3s stack is ready to get stl backend deployment
````bash
echo "install docker, k3s/k3d, go, helm, kubectl, openebs, ..." 
echo "create a cluster with a namespace : stl and registry port :5000"
cd scripts
./install_debian10.sh stl 5000

cd ~/stl
echo "change current k3s context namespace"
kubectl config set-context --current --namespace=stl
````

### Useful aliases for kubectl
For convenience, I use these aliases for running commands
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
## Install Helm mongodb
````bash
helm repo add bitnami https://charts.bitnami.com/bitnami

echo "with ClusterIp (default)"
helm install --set image.tag=3.6.23 db-stl bitnami/mongodb

echo "if you select a node port 40000 exposure"
echo 'helm install --set image.tag=3.6.23 --set service.type="NodePort" --set service.nodePort=40000 db-stl bitnami/mongodb'
````

## Install Mongodb client side (firewall, PF, noSqlClient...)
````bash
echo "MongoDB(R) can be accessed on the following DNS name(s) and ports from within your cluster:"
export MONGODB_CLUSTER_DNS="db-stl-mongodb.default.svc.cluster.local"
echo 'MONGODB_CLUSTER_DNS="db-stl-mongodb.default.svc.cluster.local"'
echo " MONGODB_CLUSTER_DNS: "$MONGODB_CLUSTER_DNS

echo "Get IP of the standalone server (CIDR > 16), first interface !"
vlan16=`ip r | grep "default" | cut -d" " -f3 | cut -d"." -f1-2`
export MONGODB_HOST=`ip -o -4 addr list | grep "$vlan16" | awk '{print $4}' | cut -d/ -f1 | head -1`
echo " IP on this vlan minimum/16: "$vlan16
echo " MONGODB_HOST: "$MONGODB_HOST

echo "To get the root password run:"
export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace stl db-stl-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)
echo 'export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace stl db-stl-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)'

echo "To connect to your database, create a MongoDB(R) client container:"
echo 'kubectl run --namespace stl db-stl-mongodb-client --rm --tty -i --restart='Never' --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:3.6.23 --command -- bash'

echo "Then, run the following command:"
echo 'mongo admin --host "db-stl-mongodb" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD'
echo 'mongo --host 127.0.0.1 --port 27017 --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD'

echo "To connect to your database from outside the cluster execute the following commands, ClusterIp :"
sleep 5
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb 27017:27017 &
echo 'kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb 27017:27017 &'
````

````sh
# if you selected NodePort deployment, port 40000 :
export NODE_IP=$(kubectl get nodes --namespace stl -o jsonpath="{.items[0].status.addresses[0].address}")
export NODE_PORT=$(kubectl get --namespace stl -o jsonpath="{.spec.ports[0].nodePort}" services mongodb-stl)
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb 27017:$NODE_PORT &
mongo --host $NODE_IP --port $NODE_PORT --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD
````
### Install noSqlclient
Launch your browser to this address http://$MONGODB_HOST:3000/
````bash
docker run -d -p 3000:3000 --name mongoclient -e MONGOCLIENT_DEFAULT_CONNECTION_URL="mongodb://root:$MONGODB_ROOT_PASSWORD@$MONGODB_HOST/admin?ssl=false" -e MONGOCLIENT_AUTH="true" -e MONGOCLIENT_USERNAME="root" -e MONGOCLIENT_PASSWORD="$MONGODB_ROOT_PASSWORD" mongoclient/mongoclient:latest
ufw allow 27017
````

### Feed with data
We need minimal data to start the application
````bash
echo "before run a mongoclient in mongo replicaset"
echo "launch a MongoShell command to install minimal data"
mgo=`kubectl get pods | grep "db-stl-mongodb" | cut -d" " -f1`
cat data/init-stl.js | sed 's/$MONGODB_ROOT_PASSWORD/'$MONGODB_ROOT_PASSWORD'/g' > /tmp/init-stl.js
kubectl cp /tmp/init-stl.js $mgo:/tmp/init-stl.js
kubectl exec -it --namespace stl $mgo -- mongo mongodb://root:$MONGODB_ROOT_PASSWORD@db-stl-mongodb:27017/ < /tmp/init-stl.js
````

### module go backend
````bash
cd backend
go get -u -v -f all
while read l; do go get -v "$l"; done < <(go list -f '{{ join .Imports "\n" }}')
go build -o stl-backend .
ufw allow 8080
./stl-backend &

````
## Development Environment
### Goland Idea
Open the project and configure these three variables to launch the rest backend 
* MONGODB_HOST
* MONGODB_USER
* MONGODB_PASSWORD

# TODO
## Development
 - [x] how to stop kubectl port-forward --namespace stl svc/db-stl-mongodb 27017:27017
 - [x] launch docker mongoclient with --name
 - [x] create db stl + collection + user
 - [x] minimal data for application running (track id=1, user id=1, ... )
 - [ ] index mongo ?
 - [ ] initial data at startup
 - [ ] get env variable in javascript mongo command, possible?
 - [x] mongodb class and pool of connections
 - [ ] load balance traefik 27017
 - [ ] travis CI/CD
 - [ ] test
 - [ ] revert back
 - [ ] change sleep(5) by loop wait while pods is actually running

## Revert Back
````sh
# port forwarding if apply
kill -9 `ps -aux | grep "27017:27017" | grep "kubectl" | awk '{print $2}'`
# revert back mongo
db.dropUser("stluser");
db.tracks.drop();
db.dropDatabase();
````
