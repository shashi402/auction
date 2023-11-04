package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/auction/logger"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Bidder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
type Bid struct {
	ID        int       `json:"id"`
	AuctionID int       `json:"auction_id"`
	BidderID  int       `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:Shashi@1234@tcp(localhost:3306)/auction")
	if err != nil {
		logger.Log.Info().Msg("Error here! ")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/bidders", listBiddersHandler)
	router.POST("/bidders", addBidderHandler)
	router.POST("/place-bid", placeBidHandler)
	log.Fatal(router.Run(":8080"))
}
func listBiddersHandler(c *gin.Context) {
	bidders, err := getBidders()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, bidders)
}

func getBidders() ([]Bidder, error) {
	rows, err := db.Query("SELECT id, name, email, phone_number FROM bidders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bidders []Bidder
	for rows.Next() {
		var bidder Bidder
		err := rows.Scan(&bidder.ID, &bidder.Name, &bidder.Email, &bidder.PhoneNumber)
		if err != nil {
			return nil, err
		}
		bidders = append(bidders, bidder)
	}

	return bidders, nil
}

func addBidderHandler(c *gin.Context) {
	var newBidder Bidder
	if err := c.ShouldBindJSON(&newBidder); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := addBidder(newBidder)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"id": id})
}

func addBidder(bidder Bidder) (int, error) {
	result, err := db.Exec("INSERT INTO bidders (name, email, phone_number) VALUES (?, ?, ?)", bidder.Name, bidder.Email, bidder.PhoneNumber)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
func placeBidHandler(c *gin.Context) {
	var newBid Bid
	if err := c.ShouldBindJSON(&newBid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := placeBid(newBid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func placeBid(bid Bid) (int, error) {
	result, err := db.Exec("INSERT INTO bids (auction_id, bidder_id, amount, timestamp) VALUES (?, ?, ?, ?)", bid.AuctionID, bid.BidderID, bid.Amount, bid.Timestamp)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
