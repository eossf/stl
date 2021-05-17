
# STL : Saintelyon Application

## Installation from ZERO - Debian 10
````bash
apt update
apt install git
git clone https://github.com/eossf/common.git
cd common
find . -name "*.sh" -exec chmod 775  {} \;
````

### create namespace stl and registry on port 5000
````bash
./install_debian10.sh stl 5000
````
### usefull aliases
````bash
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
````

### Install Helm mongodb from bitnami
````bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/mongodb
````
````bash
# (...)

MongoDB(R) can be accessed on the following DNS name(s) and ports from within your cluster: my-release-mongodb.default.svc.cluster.local

To get the root password run:
    export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace default my-release-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

To connect to your database, create a MongoDB(R) client container:
    kubectl run --namespace default my-release-mongodb-client --rm --tty -i --restart='Never' --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:4.4.6-debian-10-r0 --command -- bash

Then, run the following command:
    mongo admin --host "my-release-mongodb" --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD

To connect to your database from outside the cluster execute the following commands:
    kubectl port-forward --namespace default svc/my-release-mongodb 27017:27017 &
    mongo --host 127.0.0.1 --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD
````
