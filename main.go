package main

import (
	"log"

	"github.com/kataras/iris"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/icalF/openshop/web/controllers"
	"github.com/icalF/openshop/datasource"
	"github.com/icalF/openshop/services"
	"github.com/icalF/openshop/dao"
	"github.com/kataras/iris/_examples/mvc/overview/web/middleware"
)

func main() {
	app := iris.New()

	dbConn, err := datasource.NewMysqlConnection()
	if err != nil {
		log.Fatalf(err.Error())
	}

	userDAO := dao.NewUserDAO(dbConn)
	couponDAO := dao.NewCouponDAO(dbConn)
	productDAO := dao.NewProductDAO(dbConn)
	paymentDAO := dao.NewPaymentDAO(dbConn)
	orderDetailDAO := dao.NewOrderDetailDAO(dbConn)
	orderDAO := dao.NewOrderDAO(dbConn)

	userService := services.NewUserService(userDAO)
	productService := services.NewProductService(productDAO)
	paymentService := services.NewPaymentService(paymentDAO)
	orderDetailService := services.NewOrderDetailService(orderDetailDAO)
	orderService := services.NewOrderService(orderDAO)
	couponService := services.NewCouponService(couponDAO, orderService)

	app.Party("/", middleware.BasicAuth)

	app.Controller("/user", new(controllers.UserController),
		userService,
		// middleware.BasicAuth
	)

	app.Controller("/coupon", new(controllers.CouponController),
		couponService,
	)

	app.Controller("/product", new(controllers.ProductController),
		productService,
	)

	app.Controller("/order", new(controllers.OrderController),
		couponService,
		orderService,
		orderDetailService,
		paymentService,
	)

	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutVersionChecker,
		//iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations, // enables faster json serialization and more
	)

	defer dbConn.Close()
}
