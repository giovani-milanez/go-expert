package auction

import (
	"context"
	"fmt"
	"giovani-milanez/go-expert/auction/configuration/logger"
	"giovani-milanez/go-expert/auction/internal/entity/auction_entity"
	common "giovani-milanez/go-expert/auction/internal/infra/database"
	"giovani-milanez/go-expert/auction/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
	auctionInterval       time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
		auctionInterval: common.GetAuctionInterval(),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// vai rodar na sua propria goroutine
	_ = time.AfterFunc(ar.auctionInterval, func() {
		filter := bson.M{"_id": auctionEntity.Id}
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
		_, err := ar.Collection.UpdateOne(ctx, filter, update)	
		if err != nil {
			logger.Error("Error trying to update auction status", err)
			return
		}
		logger.Info(fmt.Sprintf("Auction status updated to completed: %s", auctionEntity.Id))
	})

	return nil
}
