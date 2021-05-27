// use stl
db = db.getSiblingDB("stl");
db.dropUser("stluser");
db.tracks.drop();
db.dropDatabase();