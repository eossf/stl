# STL : SainTeLyon flutter app
STL is a demo of a fullstack application, features are :
- [x] Live Run, follows an existing track (GPX)
- [ ] Live Run, records a track (GPX)
- [x] Support a runner: go to stops using GPS mobile
- [ ] Live View runners and position in a map
- [ ] Guide a runner
- [ ] Misc: Import, Export GPX, Enrich with metadata, ...

Connect to your server
````sh
ssh -i ~/.ssh/id_rsa root@REMOTE_IP
````
Clone this project on a remote server
````sh
apt -y install git
git clone https://github.com/eossf/stl.git
````
## Automatic installation
This script read "bash" code blocks in this file. "sh" blocks are ignored. 
````sh
cd stl
scripts/mbe.sh readme.md
````
## Manual installation
### Common tools
````bash
echo " -----------------------------------------"
echo " --- ### Common tools"
echo " -----------------------------------------"

git clone https://github.com/eossf/common.git
echo 'Done this on the git folder: find . -name "*.sh" -exec git add --chmod=+x {} \;'

````
### Full k3s/k3d installation
After this step the kubernetes k3s stack is ready to get stl backend deployment
````bash
echo " -----------------------------------------"
echo " --- ### Full k3s/k3d installation"
echo " -----------------------------------------"

echo "install docker, k3s/k3d, go, helm, kubectl, openebs, ..." 
echo "create a cluster with a namespace : stl and registry port :5000"
cd ~/stl/common/scripts
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
### Install Helm mongodb
````bash
echo " -----------------------------------------"
echo " --- ### Install Helm mongodb"
echo " -----------------------------------------"

helm repo add bitnami https://charts.bitnami.com/bitnami

echo "with ClusterIp (default)"
helm install --set image.tag=3.6.23 db-stl bitnami/mongodb

echo "if you select a node port 40000 exposure"
echo 'helm install --set image.tag=3.6.23 --set service.type="NodePort" --set service.nodePort=40000 db-stl bitnami/mongodb'
````

### Set up firewall, PF, noSqlClient...
````bash
echo " -----------------------------------------"
echo " --- ### Set up firewall, PF, noSqlClient"
echo " -----------------------------------------"

export PORT_MONGODB=27017
export MONGODB_CLUSTER_DNS="db-stl-mongodb.default.svc.cluster.local"

echo "MongoDB(R) can be accessed on the following DNS name(s) and ports from within your cluster:"$MONGODB_CLUSTER_DNS

vlan16=`ip r | grep "default" | cut -d" " -f3 | cut -d"." -f1-2`
export MONGODB_HOST=`ip -o -4 addr list | grep "$vlan16" | awk '{print $4}' | cut -d/ -f1 | head -1`
echo "IP for this vlan "$vlan16"/16 MONGODB_HOST: "$MONGODB_HOST

export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace stl db-stl-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

echo "To connect to your database, create a MongoDB(R) client container:"
echo '  kubectl run --namespace stl db-stl-mongodb-client --rm --tty -i --restart="Never" --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:3.6.23 --command -- bash'
echo "Then, run the following command:"
echo '  mongo admin --host "db-stl-mongodb" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD'
echo '  mongo --host 127.0.0.1 --port $PORT_MONGODB --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD'
echo "To connect to your database from outside the cluster execute the following commands, ClusterIp :"
echo '  kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb $PORT_MONGODB:$PORT_MONGODB &'
````

````sh
# if you selected NodePort deployment, port 40000 :
export NODE_IP=$(kubectl get nodes --namespace stl -o jsonpath="{.items[0].status.addresses[0].address}")
export NODE_PORT=$(kubectl get --namespace stl -o jsonpath="{.spec.ports[0].nodePort}" services mongodb-stl)
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb $PORT_MONGODB:$NODE_PORT &
mongo --host $NODE_IP --port $NODE_PORT --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD
````
### Install noSqlclient
Launch your browser to this address http://$MONGODB_HOST:3000/
````bash
echo " -----------------------------------------"
echo " --- ### Install noSqlclient"
echo " -----------------------------------------"

docker run -d -p 3000:3000 --name mongoclient -e MONGOCLIENT_DEFAULT_CONNECTION_URL="mongodb://root:$MONGODB_ROOT_PASSWORD@$MONGODB_HOST/admin?ssl=false" -e MONGOCLIENT_AUTH="true" -e MONGOCLIENT_USERNAME="root" -e MONGOCLIENT_PASSWORD="$MONGODB_ROOT_PASSWORD" mongoclient/mongoclient:latest
ufw allow $PORT_MONGODB
````

### Feed with configuration data
We need minimal data to start the application
````bash
echo " -----------------------------------------"
echo " --- ### Feed with configuration data"
echo " -----------------------------------------"
echo ""
echo -n "Wait after mongo pod is actually creating "
while [[ `kubectl get pods -A | grep "db-stl-mongodb" | wc -l` -eq 0 ]]; do echo -n "."; sleep 1; done
PODMONGO=`kubectl get pods | grep "db-stl-mongodb" | cut -d" " -f1`
while [[ `kubectl get pods $PODMONGO | grep "Running" | wc -l` -eq 0 ]]; do echo -n "."; sleep 1; done
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb $PORT_MONGODB:$PORT_MONGODB &
cat data/init-stl.js | sed 's/$MONGODB_ROOT_PASSWORD/'$MONGODB_ROOT_PASSWORD'/g' > /tmp/init-stl.js
kubectl exec -i --namespace stl $PODMONGO -- mongo mongodb://root:$MONGODB_ROOT_PASSWORD@127.0.0.1:$PORT_MONGODB/ < /tmp/init-stl.js
````

### STL backend
````bash
echo " -----------------------------------------"
echo " --- ### STL backend"
echo " -----------------------------------------"

cd ~/stl/backend
export PORT_STL_BACKEND=8080
go get -u -v -f all
while read l; do go get -v "$l"; done < <(go list -f '{{ join .Imports "\n" }}')
go build -o stl-backend .
ufw allow $PORT_STL_BACKEND
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
 - [ ] initial data at startup (before the container is started)
 - [ ] get env variable in javascript mongo command, possible?
 - [x] mongodb class and pool of connections
 - [ ] load balance traefik 27017
 - [ ] travis CI/CD
 - [ ] test
 - [x] revert
 - [ ] change sleep(5) by loop wait while pods is actually running

## Revert Back
````sh
cd ~/stl

export PORT_STL_BACKEND=8080
export PORT_MONGODB=27017
vlan16=`ip r | grep "default" | cut -d" " -f3 | cut -d"." -f1-2`
export MONGODB_HOST=`ip -o -4 addr list | grep "$vlan16" | awk '{print $4}' | cut -d/ -f1 | head -1`
export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace stl db-stl-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

echo " [x] revert stl-backend ----------------- "
while [[ `ufw status numbered | grep "$PORT_STL_BACKEND" | wc -l` -gt 0 ]]; do echo "delete mongo rule "$PORT_STL_BACKEND ; ufw --force delete `ufw status numbered | grep "$PORT_STL_BACKEND" | tail -1 | cut -d"[" -f2 | cut -d"]" -f1`; done
PID=`ps -a | grep "stl-backend" | cut -d" " -f1`
kill -9 $PID

echo " [x] revert data ----------------- "
kubectl exec -i --namespace stl $PODMONGO -- mongo mongodb://root:$MONGODB_ROOT_PASSWORD@127.0.0.1:$PORT_MONGODB/ < data/destroy-stl.js

echo " [x] revert noSqlclient ----------------- "
docker stop mongoclient
docker rm mongoclient

echo " [x] revert mongo ----------------- "
# port forwarding if apply
PID=`ps -aux | grep "$PORT_MONGODB:$PORT_MONGODB" | grep "kubectl" | awk '{print $2}'`
kill -9 $PID
while [[ `ufw status numbered | grep "$PORT_MONGODB" | wc -l` -gt 0 ]]; do echo "delete mongo rule "$PORT_MONGODB ; ufw --force delete `ufw status numbered | grep "$PORT_MONGODB" | tail -1 | cut -d"[" -f2 | cut -d"]" -f1`; done
helm uninstall db-stl

echo " [x] revert k3d ----------------- "
cd ~/stl/common/scripts
./delete_k3d.sh stl

echo " [x] remove project ----------------- "
cd ~
rm -Rf stl

echo " [x] restart docker ----------------- "
systemctl restart docker 
````
