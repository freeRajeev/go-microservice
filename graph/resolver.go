package graph

import (
	"context"

	"github.com/freeRajeev/go-microservice/graph/generated"
	"github.com/freeRajeev/go-microservice/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resolver struct {
	DB *mongo.Database
}

func (r *mutationResolver) CreateCustomer(ctx context.Context, name string, email string) (*model.Customer, error) {
	collection := r.DB.Collection("customers")
	customer := model.Customer{
		ID:    primitive.NewObjectID().Hex(),
		Name:  name,
		Email: email,
	}
	_, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
func (r *queryResolver) Customers(ctx context.Context) ([]*model.Customer, error) {
	collection := r.DB.Collection("customers")
	var customers []*model.Customer
	cursor, err := collection.Find(ctx, bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var customer model.Customer
		if err = cursor.Decode(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
