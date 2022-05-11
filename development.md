# STL application
## Development Environment
### Goland Idea
Open the project and configure these three variables to launch the rest backend
MONGODB_HOST=db-stl-mongodb ;MONGODB_ROOT_PASSWORD=secret;PORT_STL_BACKEND=8080;

### MongoDB local
    docker run -it --rm mongo mongo --host localhost -u admin -p secr3t --authenticationDatabase admin stl

or

    docker-compose up -d --force-recreate


### run scripts 
docker exec mongodb bash -c 'mongo < /scripts/init-stl-unsecure.js'
docker exec mongodb bash -c 'mongo < /scripts/destroy-stl.js' 

#### does not work ...
docker exec mongodb bash -c 'mongo mongodb://root:secr3t@localhost:27017/ < /scripts/init-stl-unsecure.js'


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
- [x] change sleep(5) by loop wait while pods is actually running
