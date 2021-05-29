# STL backend application
## Revert Back
````bash
cd ~/stl

export PORT_STL_BACKEND=8080
export PORT_MONGODB=27017
export MONGODB_CLUSTER_DNS="db-stl-mongodb.default.svc.cluster.local"
export PORT_NOSQLCLIENT=3000
export PODMONGO=`kubectl get pods | grep "db-stl-mongodb" | cut -d" " -f1`
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
while [[ `ufw status numbered | grep "$PORT_NOSQLCLIENT" | wc -l` -gt 0 ]]; do echo "delete mongo rule "$PORT_NOSQLCLIENT ; ufw --force delete `ufw status numbered | grep "$PORT_NOSQLCLIENT" | tail -1 | cut -d"[" -f2 | cut -d"]" -f1`; done

echo " [x] revert mongo ----------------- "
# port forwarding if apply
PID=`ps -aux | grep "$PORT_MONGODB:$PORT_MONGODB" | grep "kubectl" | awk '{print $2}'`
kill -9 $PID
while [[ `ufw status numbered | grep "$PORT_MONGODB" | wc -l` -gt 0 ]]; do echo "delete mongo rule "$PORT_MONGODB ; ufw --force delete `ufw status numbered | grep "$PORT_MONGODB" | tail -1 | cut -d"[" -f2 | cut -d"]" -f1`; done
helm uninstall db-stl

echo " [x] revert k3d ----------------- "
cd ~/stl/common/k3d
./delete_k3d.sh stl

echo " [x] remove project ----------------- "
cd ~
rm -Rf stl

echo " [x] restart docker ----------------- "
systemctl restart docker 
````
