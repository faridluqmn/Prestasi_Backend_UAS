package repository

import (
	"context"
	"prestasi_backend/app/model"
	"prestasi_backend/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAchievementColl() (*mongo.Collection, error) {
	if database.MongoDB == nil {
		return nil, errors.New("MongoDB belum terkoneksi — pastikan ConnectMongo() sudah dipanggil")
	}
	return database.MongoDB.Collection("achievements"), nil
}

// Create achievement di Mongo
func CreateAchievement(doc *model.AchievementMongo) (string, error) {
	coll, err := getAchievementColl()
	if err != nil {
		return "", err
	}

	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	oid := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

// Get achievement by ID
func GetAchievementByID(id string) (*model.AchievementMongo, error) {
	coll, err := getAchievementColl()
	if err != nil {
		return nil, err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc model.AchievementMongo
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// List achievement by studentId
func GetAchievementsByStudentID(studentID string) ([]model.AchievementMongo, error) {
	coll, err := getAchievementColl()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := coll.Find(ctx, bson.M{"studentId": studentID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.AchievementMongo
	for cur.Next(ctx) {
		var doc model.AchievementMongo
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		list = append(list, doc)
	}
	return list, cur.Err()
}

// Update achievement
func UpdateAchievement(id string, update bson.M) error {
	coll, err := getAchievementColl()
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update["updatedAt"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
	return err
}

// Delete achievement
func DeleteAchievement(id string) error {
	coll, err := getAchievementColl()
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// Tambah attachment (nama file disimpan di array attachments)
func AddAchievementAttachment(id, filename string) error {
	coll, err := getAchievementColl()
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$push": bson.M{"attachments": filename}},
	)
	return err
}
