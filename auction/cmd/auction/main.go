package main

import (
	"context"
	"giovani-milanez/go-expert/auction/configuration/database/mongodb"
	"giovani-milanez/go-expert/auction/internal/infra/api/web/controller/auction_controller"
	"giovani-milanez/go-expert/auction/internal/infra/api/web/controller/bid_controller"
	"giovani-milanez/go-expert/auction/internal/infra/api/web/controller/user_controller"
	"giovani-milanez/go-expert/auction/internal/infra/database/auction"
	"giovani-milanez/go-expert/auction/internal/infra/database/bid"
	"giovani-milanez/go-expert/auction/internal/infra/database/user"
	"giovani-milanez/go-expert/auction/internal/usecase/auction_usecase"
	"giovani-milanez/go-expert/auction/internal/usecase/bid_usecase"
	"giovani-milanez/go-expert/auction/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	userController, bidController, auctionsController := initDependencies(databaseConnection)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
