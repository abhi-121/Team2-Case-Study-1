package main

import (
	CustomerModels "Team2CaseStudy1/pkg/Customer/Models"
	OrderModels "Team2CaseStudy1/pkg/Order/Models"
	"Team2CaseStudy1/pkg/OrderProto/orderpb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var queryServiceClient orderpb.QueryServiceClient

func GetIndex(c *gin.Context) {
	// apiHitcount.Inc()
	c.JSON(http.StatusOK, gin.H{
		"Team 2": "Hello from Aadithya, Abhishek, Priya, Shashi!",
	})
}

// To get all the orders
func GetAllOrders(c *gin.Context) {
	req := &orderpb.NoParamRequest{}
	res, err := queryServiceClient.GetOrders(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get all the customers
func GetAllCustomers(c *gin.Context) {
	req := &orderpb.NoParamRequest{}
	res, err := queryServiceClient.GetCustomers(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get all the restaurants
func GetAllRestaurants(c *gin.Context) {
	req := &orderpb.NoParamRequest{}
	res, err := queryServiceClient.GetRestaurants(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.DummyRes})

}

// To get specific customer
func GetSpecificCustomer(c *gin.Context) {
	customerid, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &orderpb.SpecificCustomerRequest{CustId: customerid}

	res, err := queryServiceClient.GetACustomer(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To get specific order
func GetSpecificOrder(c *gin.Context) {
	orderid, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := &orderpb.SpecificOrderRequest{OrderId: orderid}

	res, err := queryServiceClient.GetAnOrder(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res.Res})

}

// To place a new order
func PostOrder(c *gin.Context) {

	body := c.Request.Body
	byteContent, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Sorry no content found: ", err.Error())
	}

	var NewOrder OrderModels.Order
	_ = json.Unmarshal(byteContent, &NewOrder)

	orderid := NewOrder.OrderId
	customerid := NewOrder.CustomerId
	restaurantid := NewOrder.RestaurantId
	itemlist := NewOrder.ItemLine
	price := NewOrder.Price
	discount := NewOrder.Discount

	var itemline []*orderpb.Item

	for i := range itemlist {
		itemline = append(itemline, &orderpb.Item{
			Name:  itemlist[i].Name,
			Price: itemlist[i].Price,
		})
	}

	req := &orderpb.OrderRequest{Ord: &orderpb.Order{
		OrderId:      orderid,
		CustomerId:   customerid,
		RestaurantId: restaurantid,
		ItemLine:     itemline,
		Price:        price,
		Discount:     discount,
	}}

	res, err := queryServiceClient.AddOrder(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res.Res,
	})
}

// To add new customer
func PostCustomer(c *gin.Context) {

	body := c.Request.Body
	byteContent, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Sorry no content found: ", err.Error())
	}

	var NewCustomer CustomerModels.Customer

	_ = json.Unmarshal(byteContent, &NewCustomer)

	req := &orderpb.CustomerRequest{Cust: &orderpb.Customer{
		CustomerId: NewCustomer.CustomerId,
		Name:       NewCustomer.Name,
		Address:    NewCustomer.Address,
		Phone:      NewCustomer.Phone,
	}}

	res, err := queryServiceClient.AddCustomer(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res.Res,
	})
}

// To place a new order
func PostRestaurant(c *gin.Context) {
	req := &orderpb.RestaurantRequest{}

	res, err := queryServiceClient.AddRestaurant(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": res.Res,
	})
}

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_in_celsius",
		Help: "Current teamperature of CPU in degree celsius",
	})
	apiHitcount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_api_hit_count",
		Help: "Number of times APIs were hitted",
	})
)

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(apiHitcount)
}

func apiHitCountTracker(c *gin.Context) {
	apiHitcount.Inc()
	c.Next()
}

func main() {

	cpuTemp.Set(65.3)
	fmt.Println("hello from API INIT function")
	conn, err := grpc.Dial("192.168.0.3:5051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Sorry client cannot talk to server: %v", err)
	}

	defer conn.Close()

	queryServiceClient = orderpb.NewQueryServiceClient(conn)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	apiRouter := router.Group("/api")
	apiRouter.Use(apiHitCountTracker)

	apiRouter.GET("/", GetIndex)
	apiRouter.GET("/orders", GetAllOrders)
	apiRouter.GET("/customers", GetAllCustomers)
	apiRouter.GET("/restaurants", GetAllRestaurants)
	apiRouter.GET("/order/:id", GetSpecificOrder)
	apiRouter.GET("/customer/:id", GetSpecificCustomer)

	apiRouter.POST("/new-order", PostOrder)
	apiRouter.POST("/new-customer", PostCustomer)
	apiRouter.POST("/new-restaurant", PostRestaurant)

	router.Run("0.0.0.0:9001")
}
