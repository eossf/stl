// use stl
db = db.getSiblingDB("stl");
db.createUser({
    user: "stluser",
    pwd: "$MONGODB_ROOT_PASSWORD",
    roles: [{role: "readWrite", db: "stl"}, { role: "dbAdmin", db: "stl" } ]
});
db.createCollection("tracks");
db.tracks.insertOne({"id": 1, "name": "SainteLyon 2021 JIT 75kms", "author": "stluser", steps: 1});
db.tracks.insertOne({"id": 2, "name": "UTMB 2021 CCC 101kms", "author": "stluser", steps: 1});