package main

import (
	"context"
	"giovani-milanez/go-expert/auction/configuration/database/mongodb"
	"giovani-milanez/go-expert/auction/internal/entity/auction_entity"
	"giovani-milanez/go-expert/auction/internal/infra/database/auction"
	common "giovani-milanez/go-expert/auction/internal/infra/database"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)
func TestAutomaticallyCloseAuction(t *testing.T) {
	ctx := context.Background()

	if err := godotenv.Load(".envtest"); err != nil {
		t.Fatalf("Error trying to load env variables")
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	repo := auction.NewAuctionRepository(databaseConnection)

	id := uuid.New().String()

	internal_err := repo.CreateAuction(ctx, &auction_entity.Auction{
		Id:          id,
		ProductName: "test-product-name",
		Category:    "test-category",
		Description: "test-description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	})

	if internal_err != nil {
		t.Fatalf("Error creating auction: %v", internal_err)
	}

	itvl := common.GetAuctionInterval()
	time.Sleep(itvl) // wait for the auction to close
	time.Sleep(1 * time.Second) // give some time for the goroutine to finish
	auction, internal_err := repo.FindAuctionById(ctx, id)
	if internal_err != nil {
		t.Fatalf("Error finding auction: %v", internal_err)
	}
	if auction.Status != auction_entity.Completed {
		t.Errorf("Expected auction status to be completed, got %v", auction.Status)
	}

}