package main


import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"order/orderProto"
)



type server struct{}

func (*server) CreateOrder(ctx context.Context, req *orderProto.OrderRequest) (*orderProto.OrderResponse, error){
	fmt.Println("Function called... ")
	//var obj *orderProto.OrderRequest
	orderID := req.GetOrderID()
	customerID := req.GetCustomerID()
	restaurant := req.GetRestaurant()
	//itemLine := req.GetItemLine()
	price := req.GetPrice()
	quantity := req.GetQuantity()
	discount := req.GetDiscount()
	date := req.GetDate()

	res := &orderProto.OrderResponse{
		OrderID:    orderID,
		CustomerID: customerID,
		Restaurant: restaurant,
		//ItemLine:   itemLine,
		Price:      price,
		Quantity:   quantity,
		Discount:   discount,
		Date:       date,
	}
	return res,nil

}


func (*server) GetAvgPricesOrders( ctx context.Context, req *orderProto.AvgPriceInfoRequest) (*orderProto.AvgPriceInfoResponse, error) {
	fmt.Println("AvgPrice Function called... ")
	//var obj *orderProto.OrderRequest
	customerID := req.GetCustomerID()
	avgPrice := req.GetAvgPrice()
	avgOrders := req.GetAvgOrders()

	res := &orderProto.AvgPriceInfoResponse{
		CustomerID: customerID,
		AvgPrice: float32(avgPrice),
		AvgOrders: avgOrders,
	}
	return res,nil
}


func (*server) GetTopCustomers( ctx context.Context, req *orderProto.TopCustomersRequest) (*orderProto.TopCustomersResponse, error) {
	fmt.Println("TopCustomer Function called... ")
	//var obj *orderProto.OrderRequest
	customerID := req.GetCustomerID()
	expenditure := req.GetExpenditure()

	res := &orderProto.TopCustomersResponse{
		CustomerID: customerID,
		Expenditure: float32(expenditure),
	}
	return res,nil
}

func (*server) GetTopRest( ctx context.Context, req *orderProto.TopRestaurantsRequest) (*orderProto.TopRestaurantsResponse, error) {
	fmt.Println("TopRest Function called... ")
	//var obj *orderProto.OrderRequest
	rest := req.GetRestaurant()
	revenue := req.GetRevenue()

	res := &orderProto.TopRestaurantsResponse{
		Restaurant: rest,
		Revenue: float32(revenue),
	}
	return res,nil
}


func main() {
	fmt.Println("Hello from grpc server.")

	lis, err := net.Listen("tcp","0.0.0.0:5051")
	if err!=nil {
		log.Fatalf("Sorry failed to load server %v:",err)
	}

	s := grpc.NewServer()

	orderProto.RegisterOrderServer(s, &server{})


	if s.Serve(lis); err!=nil {
		log.Fatalf("Failed to serve %v",err)
	}

}