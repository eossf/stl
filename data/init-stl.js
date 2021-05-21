// use stl
db = db.getSiblingDB("stl");
db.createUser({
    user: "MONGODB_USER",
    pwd: "MONGODB_PASSWORD",
    roles: [{role: "readWrite", db: "stl"}, { role: "dbAdmin", db: "stl" } ]
});
db.createCollection("track");
db.track.insertOne({"id": "1", "name": "Default track ID 1", "author": "stluser", steps: 1});