# Summit Tracks List


STL is a demo of a fullstack application. 

Features are :

- [ ] v.0.0.x **Alpha version:** display tracks, display menu, test persistance, gps, googlemaps, metadata, ...
- [ ] v0.1.x **View Tracks:** display tracks, current location for guiding a runner
- [ ] v0.2.x **Misc:** Import, Export GPX, Enrich with metadata, ...
- [ ] v0.3.x **Follow a runner:** find location of one or more runner, go to each steps using mobile GPS direction
- [ ] v0.4.x **Live Run:** follow an existing track, record a new track, while your runs



## Release

Current v 0.0.1



![](C:\DevOps\stl\stl\images\stl-v0.0.1.PNG)

## Installation

With an ubuntu existing server.

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
### export variables
````bash
echo " -----------------------------------------"
echo " --- ### Export variables"
echo " -----------------------------------------"

export PORT_STL_BACKEND=8080
export PORT_MONGODB=27017
export MONGODB_CLUSTER_DNS="db-stl-mongodb.default.svc.cluster.local"
export PORT_NOSQLCLIENT=3000
````
### Install Helm mongodb
````bash
echo " -----------------------------------------"
echo " --- ### Install Helm mongodb"
echo " -----------------------------------------"

ufw allow $PORT_MONGODB

helm repo add bitnami https://charts.bitnami.com/bitnami

echo "with ClusterIp (default)"
helm install db-stl bitnami/mongodb
sleep 5

echo ""
echo -n "Wait after mongo pod is actually creating "
while [[ `kubectl get pods -A | grep "db-stl-mongodb" | wc -l` -eq 0 ]]; do echo -n ":"; sleep 1; done
export PODMONGO=`kubectl get pods | grep "db-stl-mongodb" | cut -d" " -f1`
while [[ `kubectl get pods $PODMONGO | grep "Running" | wc -l` -eq 0 ]]; do echo -n "."; sleep 1; done
echo ""
````

If you select a node port 40000 exposure
````sh
helm install --set service.type="NodePort" --set service.nodePort=40000 db-stl bitnami/mongodb
````

### Set Mongo variables
````bash
echo " -----------------------------------------"
echo " --- ### Set Mongo variables"
echo " -----------------------------------------"

vlan16=`ip r | grep "default" | cut -d" " -f3 | cut -d"." -f1-4`
export MONGODB_HOST=`ip -o -4 addr list | grep "$vlan16" | awk '{print $4}' | cut -d/ -f1 | head -1`
export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace stl db-stl-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

echo "MONGODB_HOST=$MONGODB_HOST;MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD;PORT_STL_BACKEND=$PORT_STL_BACKEND"
````

To connect to your database, create a MongoDB(R) client container
````sh
kubectl run --namespace stl db-stl-mongodb-client --rm --tty -i --restart="Never" --env="PORT_MONGODB=$PORT_MONGODB" --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb --command -- bash
````
Then, run the following command:
````sh
mongo admin --host "db-stl-mongodb" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD
mongo --host "db-stl-mongodb" --port $PORT_MONGODB --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD
````
To connect to your database from outside the cluster execute the following commands, ClusterIp :
 ````sh
kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb $PORT_MONGODB:$PORT_MONGODB &
 ````
If you selected NodePort deployment, port 40000 :
 ````sh
 export NODE_IP=$(kubectl get nodes --namespace stl -o jsonpath="{.items[0].status.addresses[0].address}")
 export NODE_PORT=$(kubectl get --namespace stl -o jsonpath="{.spec.ports[0].nodePort}" services mongodb-stl)
 kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb $PORT_MONGODB:$NODE_PORT &
 mongo --host $NODE_IP --port $NODE_PORT --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD
 ````

````bash
echo " -----------------------------------------"
echo " --- ### Port forwarding "
echo " -----------------------------------------"

kubectl port-forward --namespace stl --address 0.0.0.0 svc/db-stl-mongodb $PORT_MONGODB:$PORT_MONGODB &
sleep 5
````
### Install noSqlclient

````bash
echo " -----------------------------------------"
echo " --- ### Install noSqlclient"
echo " -----------------------------------------"

docker run -d -p $PORT_NOSQLCLIENT:$PORT_NOSQLCLIENT --name mongoclient -e MONGOCLIENT_DEFAULT_CONNECTION_URL="mongodb://root:$MONGODB_ROOT_PASSWORD@$MONGODB_HOST/admin?ssl=false" -e MONGOCLIENT_AUTH="true" -e MONGOCLIENT_USERNAME="root" -e MONGOCLIENT_PASSWORD="$MONGODB_ROOT_PASSWORD" mongoclient/mongoclient:latest
ufw allow $PORT_NOSQLCLIENT
````

Launch your browser to this address http://$MONGODB_HOST:3000/

### Feed with configuration data
We need minimal data to start the application
````bash
echo " -----------------------------------------"
echo " --- ### Feed with configuration data"
echo " -----------------------------------------"
echo ""

cat data/init-stl.js | sed 's/$MONGODB_ROOT_PASSWORD/'$MONGODB_ROOT_PASSWORD'/g' > /tmp/init-stl.js
kubectl exec -i --namespace stl $PODMONGO -- mongo mongodb://root:$MONGODB_ROOT_PASSWORD@$MONGODB_HOST:$PORT_MONGODB/ < /tmp/init-stl.js
````

### STL backend
````bash
echo " -----------------------------------------"
echo " --- ### STL backend"
echo " -----------------------------------------"

cd ~/stl/
go mod tidy
cd ~/stl/backend
go get -u -v -f all
while read l; do go get -v "$l"; done < <(go list -f '{{ join .Imports "\n" }}')
go build -o stl-backend .
ufw allow $PORT_STL_BACKEND
./stl-backend &
````