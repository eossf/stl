
db.createUser({
    user: 'admin',
    pwd: 'secret',
    roles: [{role: "readWrite", db: "admin"}, {role: "dbAdmin", db: "admin"}]
});

db = new Mongo().getDB("stl");
db = db.getSiblingDB("stl");
db.createUser({
    user: "stluser",
    pwd: "stluser",
    roles: [{role: "readWrite", db: "stl"}, {role: "dbAdmin", db: "stl"}]
});
db.createCollection("tracks", {capped: false});
db.tracks.insertOne({"id": 1, "name": "Track Test 1 - 1km", "author": "stluser", steps: 0, "starting-point": "46.1844162,6.2699820"});
db.tracks.insertOne({"id": 2, "name": "Track Test 2 - 4.5kms", "author": "stluser", steps: 1, "starting-point": "46.1844162,6.2699820"});
db.tracks.insertOne({"id": 3, "name": "Saintelyon 2019 - 76kms", "author": "stluser", steps: 10, "starting-point": "45.4504121,4.3984407"});
db.tracks.insertOne({"id": 4, "name": "Saintelyon 2021 - JIT 75kms", "author": "stluser", steps: 10, "starting-point": "45.4504121,4.3984407"});
db.tracks.insertOne({"id": 5, "name": "Saintelyon 2022 - 78kms", "author": "stluser", steps: 10, "starting-point": "45.4504121,4.3984407"});
db.tracks.insertOne({"id": 6, "name": "Montée Sur Cou - 15kms", "author": "stluser", steps: 3, "starting-point": "46.0235016,6.3233876"});