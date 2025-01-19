package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Maja-Service/internal/user/config"
	"github.com/hoitek/Maja-Service/internal/user/domain"
	"github.com/hoitek/Maja-Service/internal/user/models"
	"github.com/hoitek/Maja-Service/internal/user/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type UserRepositoryMongoDB struct {
	MongoDB *mongo.Client
}

func NewUserRepositoryMongoDB(d *mongo.Client) *UserRepositoryMongoDB {
	return &UserRepositoryMongoDB{
		MongoDB: d,
	}
}

func (r *UserRepositoryMongoDB) Query(queries *models.UsersQueryRequestParams) ([]*domain.User, error) {
	if queries == nil {
		return []*domain.User{}, nil
	}

	// Filters
	filter := bson.M{}
	f, err := utils.Jsonify(queries.Filters)
	if err != nil {
		return []*domain.User{}, err
	}
	fMap, ok := f.(map[string]interface{})
	if !ok {
		return []*domain.User{}, errors.New("failed to parse filters")
	}
	for k, v := range fMap {
		switch value := v.(type) {
		case map[string]interface{}:
			op := fmt.Sprintf("%v", value["op"])
			val := fmt.Sprintf("%v", value["value"])
			switch op {
			case operators.EQUALS:
				filter[k] = val
			case operators.NUMBER_NOT_EQUALS:
				filter[k] = bson.M{"$ne": val}
			case operators.NUMBER_GREATER_THAN:
				filter[k] = bson.M{"$gt": val}
			case operators.NUMBER_GREATER_THAN_EQUALS:
				filter[k] = bson.M{"$gte": val}
			case operators.NUMBER_LESS_THAN:
				filter[k] = bson.M{"$lt": val}
			case operators.NUMBER_LESS_THAN_EQUALS:
				filter[k] = bson.M{"$lte": val}
			case operators.CONTAINS, operators.STARTS_WITH, operators.ENDS_WITH:
				filter[k] = bson.M{"$regex": val}
			case operators.IS_ANY_OF:
				filter[k] = bson.M{"$in": val}
			}
		default:
			return []*domain.User{}, errors.New("filter not supported")
		}
	}
	if queries.ID > 0 {
		filter["id"] = queries.ID
	}
	// Pagination
	findOptions := options.Find()
	if queries.Limit > 0 {
		findOptions.SetLimit(int64(queries.Limit))
	}
	var skip = int64((queries.Page - 1) * queries.Limit)
	if skip > 0 {
		findOptions.SetSkip(skip)
	}

	// Sorts
	s, err := utils.Jsonify(queries.Sorts)
	if err != nil {
		return []*domain.User{}, err
	}
	sMap, ok := s.(map[string]interface{})
	if !ok {
		return []*domain.User{}, errors.New("failed to parse sorts")
	}
	var sorts bson.D
	for k, v := range sMap {
		switch value := v.(type) {
		case map[string]interface{}:
			op := fmt.Sprintf("%v", value["op"])
			switch op {
			case "asc":
				sorts = append(sorts, bson.E{Key: k, Value: 1})
			case "desc":
				sorts = append(sorts, bson.E{Key: k, Value: -1})
			}
		}
	}
	if len(sorts) > 0 {
		findOptions.SetSort(sorts)
	}

	cur, err := r.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).Find(context.Background(), filter, findOptions)
	if err != nil {
		return []*domain.User{}, err
	}
	defer cur.Close(context.Background())

	var users []*domain.User
	for cur.Next(context.Background()) {
		var user domain.User
		err := cur.Decode(&user)
		if err != nil {
			return []*domain.User{}, err
		}
		users = append(users, &user)
	}
	if err := cur.Err(); err != nil {
		return []*domain.User{}, err
	}
	return users, nil
}

func (r *UserRepositoryMongoDB) Count(queries *models.UsersQueryRequestParams) (int64, error) {
	if queries == nil {
		return 0, nil
	}
	filter := bson.M{}
	if queries.ID > 0 {
		filter["id"] = queries.ID
	}
	if queries.Filters.Email.Value != "" {
		filter["email"] = bson.M{"$regex": queries.Filters.Email.Value}
	}
	if queries.Filters.FirstName.Value != "" {
		filter["firstName"] = bson.M{"$regex": queries.Filters.FirstName.Value}
	}
	if queries.Filters.LastName.Value != "" {
		filter["lastName"] = bson.M{"$regex": queries.Filters.LastName.Value}
	}
	if queries.Filters.Phone.Value != "" {
		filter["phone"] = bson.M{"$regex": queries.Filters.Phone.Value}
	}
	if queries.Filters.WorkPhoneNumber.Value != "" {
		filter["workPhoneNumber"] = bson.M{"$regex": queries.Filters.WorkPhoneNumber.Value}
	}
	if queries.Filters.Username.Value != "" {
		filter["username"] = bson.M{"$regex": queries.Filters.Username.Value}
	}
	if queries.Filters.AvatarUrl.Value != "" {
		filter["avatarUrl"] = bson.M{"$regex": queries.Filters.AvatarUrl.Value}
	}
	count, err := r.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).CountDocuments(context.Background(), filter)
	log.Println("filter", err, filter, count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepositoryMongoDB) Create(payload map[string]interface{}) (interface{}, error) {
	// Check if user already exists
	userId, ok := payload["id"]
	if !ok {
		return nil, errors.New("user id is required")
	}
	count, err := r.Count(&models.UsersQueryRequestParams{
		ID: int(userId.(float64)),
	})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	// Create user
	result, err := r.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).InsertOne(context.Background(), payload)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *UserRepositoryMongoDB) Update(payload *domain.User, id int) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": payload}
	_, err := r.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryMongoDB) Delete(ids []uint) error {
	filter := bson.M{"id": bson.M{"$in": ids}}
	log.Println("----------------------", filter)
	_, err := r.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// UpdateByPostgresID updates doc by postgres id
func (r *UserRepositoryMongoDB) UpdateByPostgresID(postgresID int, payload interface{}) (interface{}, error) {
	// Check if doc already exists
	log.Println("postgresID", postgresID)
	count, err := r.Count(&models.UsersQueryRequestParams{
		ID: postgresID,
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("user does not exist")
	}

	// Update user
	result, err := r.MongoDB.Database(config.UserConfig.DatabaseMongoDBName).Collection(domain.NewUser().TableName()).UpdateOne(context.Background(), bson.M{"id": postgresID}, bson.M{"$set": payload})
	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}
