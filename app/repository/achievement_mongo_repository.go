package repository

import (
	"context"
	"errors"
	"time"

	"prestasi_backend/app/model"
	"prestasi_backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// =====================================
// ACHIEVEMENTS (MongoDB)
// =====================================

func getAchievementCollection() (*mongo.Collection, error) {
	if database.MongoDB == nil {
		return nil, errors.New("MongoDB belum terkoneksi – panggil ConnectMongo() di main.go")
	}
	return database.MongoDB.Collection("achievements"), nil
}

// @Summary      Submit Prestasi
// @Description  Mahasiswa menambahkan laporan prestasi baru ke MongoDB [cite: 178-185]
// @Tags         achievements
// @Security     BearerAuth
// @Param        achievement  body  model.AchievementCreateRequest  true  "Data Prestasi"
// @Success      201  {object}  model.AchievementMongo
// @Router       /achievements [post]
var CreateAchievement = func(doc *model.AchievementMongo) (string, error) {
	coll, err := getAchievementCollection()
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

// @Summary      Get Detail Prestasi
// @Description  Mengambil detail data prestasi dari MongoDB berdasarkan ID [cite: 210, 247]
// @Tags         achievements
// @Security     BearerAuth
// @Param        id   path   string  true  "Mongo Achievement ID"
// @Success      200  {object}  model.AchievementMongo
// @Router       /achievements/{id} [get]
func GetAchievementByID(id string) (*model.AchievementMongo, error) {
	coll, err := getAchievementCollection()
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

// Ambil semua achievement berdasarkan studentId
var GetAchievementsByStudentID = func(studentID string) ([]model.AchievementMongo, error) {
	coll, err := getAchievementCollection()
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

// Update achievement (title, description, dsb.)
func UpdateAchievement(id string, update bson.M) error {
	coll, err := getAchievementCollection()
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if update == nil {
		update = bson.M{}
	}
	update["updatedAt"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": update},
	)
	return err
}

// @Summary      Hapus Prestasi
// @Description  Menghapus data prestasi dari MongoDB (Hard Delete) [cite: 203]
// @Tags         achievements
// @Security     BearerAuth
// @Param        id   path   string  true  "Mongo Achievement ID"
// @Success      200  {object}  map[string]string "Success message"
// @Router       /achievements/{id} [delete]
func DeleteAchievement(id string) error {
	coll, err := getAchievementCollection()
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

// Tambah attachment (nama file) ke array attachments
func AddAchievementAttachment(id, filename string) error {
	coll, err := getAchievementCollection()
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
