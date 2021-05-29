# STL application
## Development Environment
### Goland Idea
Open the project and configure these three variables to launch the rest backend
MONGODB_HOST=;
MONGODB_ROOT_PASSWORD=;
PORT_STL_BACKEND=

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
