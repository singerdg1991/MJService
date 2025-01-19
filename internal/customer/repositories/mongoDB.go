package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/hoitek/Go-Quilder/operators"
	"github.com/hoitek/Kit/exp"
	"github.com/hoitek/Maja-Service/internal/_shared/utils"
	"github.com/hoitek/Maja-Service/internal/customer/config"
	"github.com/hoitek/Maja-Service/internal/customer/domain"
	"github.com/hoitek/Maja-Service/internal/customer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type CustomerRepositoryMongoDB struct {
	MongoDB *mongo.Client
}

func NewCustomerRepositoryMongoDB(d *mongo.Client) *CustomerRepositoryMongoDB {
	return &CustomerRepositoryMongoDB{
		MongoDB: d,
	}
}

func (r *CustomerRepositoryMongoDB) Query(queries *models.CustomersQueryRequestParams) ([]*domain.Customer, error) {
	var (
		limit = 1
		page  = 1
	)
	if queries != nil {
		limit = exp.TerIf(queries.Limit < 1, 1, queries.Limit)
		page = exp.TerIf(queries.Page < 1, 1, queries.Page)
	}

	var pipeline []bson.M

	// Filters
	filter := bson.M{}

	// Prepare filters
	if queries != nil {
		f, err := utils.Jsonify(queries.Filters)
		if err != nil {
			return []*domain.Customer{}, err
		}
		fMap, ok := f.(map[string]interface{})
		if !ok {
			return []*domain.Customer{}, errors.New("failed to parse filters")
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
				return []*domain.Customer{}, errors.New("filter not supported")
			}
		}
		if queries.ID > 0 {
			filter["id"] = queries.ID
		}
		if queries.UserID > 0 {
			filter["user.id"] = queries.UserID
		}

		// Sorts
		s, err := utils.Jsonify(queries.Sorts)
		if err != nil {
			return []*domain.Customer{}, err
		}
		sMap, ok := s.(map[string]interface{})
		if !ok {
			return []*domain.Customer{}, errors.New("failed to parse sorts")
		}
		for k, v := range sMap {
			sortField := k
			sortOrder := 1 // Default sort order is ascending
			if len(v.(map[string]interface{})) == 0 {
				continue
			}

			switch value := v.(type) {
			case map[string]interface{}:
				op := fmt.Sprintf("%v", value["op"])
				if op == "desc" {
					sortOrder = -1 // Set sort order to descend
				}
			}

			// Split the nested field path if it exists
			parts := strings.Split(sortField, ".")
			if len(parts) > 1 {
				nestedField := strings.Join(parts[:len(parts)-1], ".")
				sortField = parts[len(parts)-1]
				pipeline = append(pipeline, bson.M{"$addFields": bson.M{sortField: "$" + nestedField + "." + sortField}})
			}

			pipeline = append(pipeline, bson.M{"$sort": bson.M{sortField: sortOrder}})
		}
	}

	// Add the match stage to the pipeline
	pipeline = append(pipeline, bson.M{"$match": filter})

	// Add the pagination stages to the pipeline
	var skip = int64((page - 1) * page)
	if skip > 0 {
		pipeline = append(pipeline, bson.M{"$skip": skip})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	cur, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).Aggregate(context.Background(), pipeline)
	if err != nil {
		return []*domain.Customer{}, err
	}
	defer cur.Close(context.Background())

	var customers []*domain.Customer
	for cur.Next(context.Background()) {
		var customer domain.Customer
		err := cur.Decode(&customer)
		if err != nil {
			return []*domain.Customer{}, err
		}
		customers = append(customers, &customer)
	}
	if err := cur.Err(); err != nil {
		return []*domain.Customer{}, err
	}
	return customers, nil
}

func (r *CustomerRepositoryMongoDB) Count(queries *models.CustomersQueryRequestParams) (int64, error) {
	// Filters
	filter := bson.M{}

	// Prepare filters
	if queries != nil {
		f, err := utils.Jsonify(queries.Filters)
		if err != nil {
			return 0, err
		}
		fMap, ok := f.(map[string]interface{})
		if !ok {
			return 0, errors.New("failed to parse filters")
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
				return 0, errors.New("filter not supported")
			}
		}
		if queries.ID > 0 {
			filter["id"] = queries.ID
		}
		if queries.UserID > 0 {
			filter["user.id"] = queries.UserID
		}
	}

	count, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CustomerRepositoryMongoDB) Create(postgresID int, payload interface{}) (interface{}, error) {
	// Check if customer already exists
	count, err := r.Count(&models.CustomersQueryRequestParams{
		ID: postgresID,
	})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("customer already exists")
	}

	// Create customer
	result, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).InsertOne(context.Background(), payload)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *CustomerRepositoryMongoDB) Update(payload interface{}, id int) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": payload}
	_, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *CustomerRepositoryMongoDB) CreateOrUpdate(postgresID int, payload interface{}) (interface{}, error) {
	// Check if customer already exists
	customers, err := r.Query(&models.CustomersQueryRequestParams{
		ID: postgresID,
	})
	if err != nil {
		return nil, err
	}
	if len(customers) > 0 {
		// Update customer
		err := r.Update(payload, postgresID)
		if err != nil {
			return nil, err
		}
		return customers[0].MongoID, nil
	}

	// Create customer
	result, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).InsertOne(context.Background(), payload)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *CustomerRepositoryMongoDB) Delete(ids []uint) error {
	filter := bson.M{"id": bson.M{"$in": ids}}
	_, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// UpdateByPostgresID updates doc by postgres id
func (r *CustomerRepositoryMongoDB) UpdateByPostgresID(postgresID int, payload interface{}) (interface{}, error) {
	// Check if doc already exists
	count, err := r.Count(&models.CustomersQueryRequestParams{
		ID: postgresID,
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("customer does not exist")
	}

	// Create absences
	result, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).UpdateOne(context.Background(), bson.M{"id": postgresID}, bson.M{"$set": payload})
	if err != nil {
		return nil, err
	}
	return result.UpsertedID, nil
}

// UpdateUserInfo updates customer's user info
func (r *CustomerRepositoryMongoDB) UpdateUserInfo(userID int, payload interface{}) (interface{}, error) {
	// Check if doc already exists
	count, err := r.Count(&models.CustomersQueryRequestParams{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("customer does not exist")
	}

	// Update customer's user info

	newUserInfo := map[string]interface{}{
		"id":                 userID,
		"username":           payload.(map[string]interface{})["username"],
		"firstName":          payload.(map[string]interface{})["firstName"],
		"lastName":           payload.(map[string]interface{})["lastName"],
		"email":              payload.(map[string]interface{})["email"],
		"phone":              payload.(map[string]interface{})["phone"],
		"role":               payload.(map[string]interface{})["role"],
		"avatarUrl":          payload.(map[string]interface{})["avatarUrl"],
		"workPhoneNumber":    payload.(map[string]interface{})["workPhoneNumber"],
		"gender":             payload.(map[string]interface{})["gender"],
		"accountNumber":      payload.(map[string]interface{})["accountNumber"],
		"telephone":          payload.(map[string]interface{})["telephone"],
		"registrationNumber": payload.(map[string]interface{})["registrationNumber"],
		"languageSkills":     payload.(map[string]interface{})["languageSkills"],
	}

	filter := bson.M{"user.id": userID}
	update := bson.M{"$set": bson.M{"user": newUserInfo}}
	result, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

// CreateEmptyForUserID creates customer
func (r *CustomerRepositoryMongoDB) CreateEmptyForUserID(payload interface{}) (interface{}, error) {
	result, err := r.MongoDB.Database(config.CustomerConfig.DatabaseMongoDBName).Collection(domain.NewCustomer().TableName()).InsertOne(context.Background(), payload)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}
