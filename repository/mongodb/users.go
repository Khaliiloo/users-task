package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"users-task/configs"
	"users-task/helpers"
	"users-task/models"
)

func (db *MongoDB) List(ctx context.Context) (*[]models.User, error) {
	userCollection := db.getCollection("usersDB", "users")
	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if results == nil {
		return nil, fmt.Errorf("couldn't get all users")
	}

	defer results.Close(ctx)

	var users []models.User
	for results.Next(ctx) {
		var user models.User
		if err = results.Decode(&user); err != nil {
			configs.Logger.Errorf("%q", err)
		}
		if user.ID != 0 { // get user if is at least decoded right for `id`
			users = append(users, user)
		}
	}
	return &users, nil
}

func (db *MongoDB) Get(ctx context.Context, id int) (*models.User, error) {
	userCollection := db.getCollection("usersDB", "users")
	var user = models.User{}
	filter := bson.M{"id": id}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	return &user, err

}

func (db *MongoDB) Create(ctx context.Context, user *models.User) (int, error) {
	userCollection := db.getCollection("usersDB", "users")
	age := helpers.CalculateAge(user.DOB.Time)
	if age < 18 {
		return 0, fmt.Errorf("age of the user is %v, only 18 years old or older are allowed", age)
	}
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, fmt.Errorf("couldn't insert user with id %v", user.ID)
	}
	return user.ID, nil
}

func (db *MongoDB) Update(ctx context.Context, id int, user *models.User) (*models.User, error) {
	userCollection := db.getCollection("usersDB", "users")
	age := helpers.CalculateAge(user.DOB.Time)
	if age < 18 {
		return nil, fmt.Errorf("age of the user is %v, only 18 years old or older are allowed", age)
	}
	update := bson.M{"id": user.ID, "name": user.Name, "email": user.Email, "dob": user.DOB}
	result, err := userCollection.UpdateOne(ctx, bson.D{{"id", id}}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("document with id %v which shold be updated doesn't exist\n", id)
	}
	return user, nil
}

func (db *MongoDB) Delete(ctx context.Context, id int) error {
	userCollection := db.getCollection("usersDB", "users")
	result, err := userCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount < 1 {
		return fmt.Errorf("user with specified id %v not found", id)
	}
	return nil
}

func (db *MongoDB) AddFile(ctx context.Context, id int) error {
	return nil
}
