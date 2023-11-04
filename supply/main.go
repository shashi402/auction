package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"auction/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

// AdSpace represents an ad space.
type AdSpace struct {
	ID        int     `json:"id"`
	BasePrice float64 `json:"base_price"`
}
type Auction struct {
	ID          int            `json:"id"`
	AdSpaceID   int            `json:"ad_space_id"`
	StartingBid float64        `json:"starting_bid"`
	EndTime     mysql.NullTime `json:"end_time"`
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:Shashi@1234@tcp(localhost:3306)/auction")
	if err != nil {
		logger.Log.Info().Msg("")
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/adspaces", listAdSpacesHandler)
	router.GET("/auctions", listAuctionsHandler)
	router.POST("/publish-auction", publishAuctionHandler)
	router.Run(":8080")
}
func listAdSpacesHandler(c *gin.Context) {
	adSpaces, err := getAdSpaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adSpaces)
}
func getAdSpaces() ([]AdSpace, error) {
	rows, err := db.Query("SELECT * FROM ad_spaces")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adSpaces []AdSpace
	for rows.Next() {
		var adSpace AdSpace
		err := rows.Scan(&adSpace.ID, &adSpace.BasePrice)
		if err != nil {
			return nil, err
		}
		adSpaces = append(adSpaces, adSpace)
	}

	return adSpaces, nil
}
func publishAuctionHandler(c *gin.Context) {
	var newAuction Auction
	if err := c.ShouldBindJSON(&newAuction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := publishAuction(newAuction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func publishAuction(auction Auction) (int, error) {
	result, err := db.Exec("INSERT INTO auctions (ad_space_id, starting_bid, end_time) VALUES (?, ?, ?)", auction.AdSpaceID, auction.StartingBid, auction.EndTime)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// New handlers
func listAuctionsHandler(c *gin.Context) {
	auctions, err := getAuctions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func getAuctions() ([]Auction, error) {
	rows, err := db.Query("SELECT id, ad_space_id, starting_bid, end_time FROM auctions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []Auction
	for rows.Next() {
		var auction Auction
		err := rows.Scan(&auction.ID, &auction.AdSpaceID, &auction.StartingBid, &auction.EndTime)
		if err != nil {
			return nil, err
		}
		auctions = append(auctions, auction)
	}

	return auctions, nil
}
