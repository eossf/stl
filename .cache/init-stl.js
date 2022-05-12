// use stl
db = db.getSiblingDB("stl");
db.createUser({
    user: "stluser",
    pwd: "$MONGODB_ROOT_PASSWORD",
    roles: [{role: "readWrite", db: "stl"}, { role: "dbAdmin", db: "stl" } ]
});
db.createCollection("tracks");
db.tracks.insertOne({"id": 1, "name": "Track Test 1 - 1km", "author": "stluser", steps: 0});
db.tracks.insertOne({"id": 2, "name": "Track Test 2 - 4.5kms", "author": "stluser", steps: 1});
db.tracks.insertOne({"id": 3, "name": "Saintelyon - JIT 75kms", "author": "stluser", steps: 10});
