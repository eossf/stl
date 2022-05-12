db.createUser({
    user: 'root',
    pwd: 'secret',
    roles: [
        {
            role: 'readWrite',
            db: 'stl',
        },
    ],
});

db = new Mongo().getDB("stl");
db = db.getSiblingDB("stl");
db.createUser({
    user: "stluser",
    pwd: "stluser",
    roles: [{role: "readWrite", db: "stl"}, {role: "dbAdmin", db: "stl"}]
});
db.createCollection("tracks", {capped: false});
db.tracks.insertOne({"id": 1, "name": "Track Test 1 - 1km", "author": "stluser", steps: 0});
db.tracks.insertOne({"id": 2, "name": "Track Test 2 - 4.5kms", "author": "stluser", steps: 1});
db.tracks.insertOne({"id": 3, "name": "Saintelyon - JIT 75kms", "author": "stluser", steps: 10});
