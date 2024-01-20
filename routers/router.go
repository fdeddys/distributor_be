package routers

import (
	"distribution-system-be/constants"
	"distribution-system-be/controllers"
	"distribution-system-be/models"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/utils/security"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	kons "distribution-system-be/constants"

	"github.com/astaxie/beego"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	fmt.Println(gin.IsDebugging())
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		//AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge: 86400,
	}))

	UserController := new(controllers.UserController)
	CustomerController := new(controllers.CustomerController)
	SupplierController := new(controllers.SupplierController)
	OrderController := new(controllers.OrderController)
	OrderDetailController := new(controllers.OrderDetailController)
	DashboardController := new(controllers.DashboardController)
	ReceiveController := new(controllers.ReceiveController)
	ReceiveDetailController := new(controllers.ReceiveDetailController)

	RetusnSalesOrderController := new(controllers.RetusnSalesOrderController)
	ReturnSalesOrderDetailController := new(controllers.ReturnSalesOrderDetailController)

	RetusnReceiveController := new(controllers.ReturnReceiveController)
	ReturnReceiveDetailController := new(controllers.ReturnReceiveDetailController)

	AdjustmentController := new(controllers.AdjustmentController)
	AdjustmentDetailController := new(controllers.AdjustmentDetailController)

	WarehouseController := new(controllers.WarehouseController)
	SalesmanController := new(controllers.SalesmanController)

	PaymentController := new(controllers.PaymentController)
	PaymentDetailController := new(controllers.PaymentDetailController)
	PaymentOrderController := new(controllers.PaymentOrderController)
	PaymentReturnController := new(controllers.PaymentReturnController)

	StockMutationController := new(controllers.StockMutationController)
	StockMutationDetailController := new(controllers.StockMutationDetailController)

	StockController := new(controllers.StockController)

	HistoryController := new(controllers.HistoryStockController)

	PurchaseOrderController := new(controllers.PurchaseOrderController)
	PurchaseOrderDetailController := new(controllers.PurchaseOrderDetailController)

	DirectPaymentController := new(controllers.DirectPaymentController)

	StockOpnameController := new(controllers.StockOpnameController)
	StockOpnameDetailController := new(controllers.StockOpnameDetailController)

	ReportPaymentCashController := new(controllers.ReportPaymentCashController)
	ReportPaymentSupplierController := new(controllers.ReportPaymentSupplierController)

	ReportMasterController := new(controllers.ReportMasterController)

	PaymentSupplierController := new(controllers.PaymentSupplierController)
	PaymentSupplierDetailController := new(controllers.PaymentSupplierDetailController)

	AuthController := new(controllers.AuthController)

	ParameterController := new(controllers.ParameterController)

	api := r.Group("/api/user")
	api.POST("/filter/page/:page/count/:count", UserController.GetUser)
	api.POST("", UserController.SaveDataUser)
	api.PUT("", UserController.UpdateUser)
	api.POST("/reset/:iduser", cekToken, UserController.ResetUser)
	api.GET("/current-user", AuthController.GetCurrPass)

	api = r.Group("/auth")
	api.POST("/login", AuthController.Login)
	api.POST("/update-password", AuthController.ChangePass)

	api = r.Group("/api/customer")
	api.POST("/page/:page/count/:count", CustomerController.FilterDataCustomer)
	api.POST("", CustomerController.SaveDataCustomer)
	// api.PUT("/", CustomerController.EditDataCustomer)
	api.POST("/list", CustomerController.ListDataCustomerByName)
	// api.POST("/check/supplier", MerchantController.CheckOrderMerchantSupplier)

	api = r.Group("/api/supplier")
	api.POST("/page/:page/count/:count", SupplierController.FilterDataSupplier)
	api.GET("/id/:id", SupplierController.FilterByID)
	api.POST("", SupplierController.SaveDataSupplier)
	api.PUT("", SupplierController.EditDataSupplier)

	BrandController := new(controllers.BrandController)
	brand := r.Group("/api/brand")
	brand.POST("/page/:page/count/:count", BrandController.GetBrand)
	brand.GET("/id/:id", BrandController.GetFilterBrand)
	brand.POST("", BrandController.SaveBrand)
	brand.PUT("", BrandController.UpdateBrand)
	brand.GET("", BrandController.GetBrandLike)

	ProductController := new(controllers.ProductController)
	product := r.Group("/api/product")
	product.POST("/page/:page/count/:count", cekToken, ProductController.GetProductListPaging)
	product.POST("/all-status/page/:page/count/:count", cekToken, ProductController.GetProductListPagingAllStatus)
	product.POST("/search/page/:page/count/:count", cekToken, ProductController.SearchProduct)
	product.GET("/id/:id", ProductController.GetProductDetails)
	product.POST("", cekToken, ProductController.SaveProduct)
	product.GET("/list", ProductController.ProductList)
	product.GET("", ProductController.GetProductLike)
	product.GET("/processCSV", ProductController.ProcessCSV)
	product.GET("/process-update", ProductController.ProcessUpdateProd)

	ProductGroupController := new(controllers.ProductGroupController)
	productGroup := r.Group("/api/product-group")
	productGroup.POST("/page/:page/count/:count", ProductGroupController.GetProductGroupPaging)
	productGroup.GET("/id/:id", ProductGroupController.GetProductGroupDetails)
	productGroup.POST("", ProductGroupController.SaveProductGroup)
	productGroup.PUT("", ProductGroupController.UpdateProductGroup)

	LookupController := new(controllers.LookupController)
	lookup := r.Group("/api/lookup")
	lookup.GET("", LookupController.GetLookupByGroup)
	lookup.POST("/page/:page/count/:count", LookupController.GetLookupPaging)
	lookup.GET("/id/:id", LookupController.GetLookupFilter)
	lookup.GET("/name/:name", LookupController.GetLookupGroupName)
	lookup.GET("/group", LookupController.GetDistinctLookup)
	lookup.POST("", LookupController.SaveLookup)
	// lookup.PUT("/", LookupController.UpdateLookup)

	LookupGroupController := new(controllers.LookupGroupController)
	lookupGroup := r.Group("/api/lookup-group")
	lookupGroup.GET("", LookupGroupController.GetLookupGroup)

	RoleController := new(controllers.RoleController)
	api = r.Group("/api/role")
	api.POST("/filter/page/:page/count/:count", RoleController.GetRole)
	api.POST("", cekToken, RoleController.SaveRole)
	api.PUT("", RoleController.UpdateRole)

	AccMatrixController := new(controllers.AccessMatrixController)
	MenuController := new(controllers.MenuController)
	api = r.Group("/api/menu")
	api.GET("/list-user-menu", cekToken, MenuController.GetMenuByUser)
	api.GET("/list-all-active-menu", AccMatrixController.GetAllActiveMenu)
	api.GET("/role/:roleId", AccMatrixController.GetMenuByRoleID)
	api.POST("/role/:roleId", AccMatrixController.SaveRoleMenu)
	api.POST("/update-role", cekToken, AccMatrixController.UpdateRoleMenu)

	api = r.Group("/api/sales-order")
	api.GET("/:id", OrderController.GetByOrderId)
	api.GET("/:id/total", cekToken, OrderController.CekTotal)
	api.POST("/page/:page/count/:count", cekToken, OrderController.FilterData)
	api.POST("", cekToken, OrderController.Save)
	api.POST("/approve", cekToken, OrderController.Approve)
	api.POST("/create-invoice", cekToken, OrderController.CreateInvoice)
	api.POST("/reject", cekToken, OrderController.Reject)
	api.POST("/print/so/:id", OrderController.PrintSO)
	api.POST("/print/invoice/:id", OrderController.PrintInvoice)
	api.POST("/payment/page/:page/count/:count", cekToken, OrderController.FilterDataForSalesOrder)

	api = r.Group("/api/sales-order-detail")
	api.POST("/page/:page/count/:count", OrderDetailController.GetDetail)
	api.POST("", cekToken, OrderDetailController.Save)
	api.POST("/updateItemRecv", cekToken, OrderDetailController.UpdateQtyReceive)
	api.POST("/updateQtyOrder", cekToken, OrderDetailController.UpdateQtyOrder)
	api.DELETE("/:id", cekToken, OrderDetailController.DeleteById)

	// RETURN SALES-ORDER
	api = r.Group("/api/return-sales-order")
	api.GET("/:id", cekToken, RetusnSalesOrderController.GetByReturnSalesOrderId)
	api.POST("/page/:page/count/:count", cekToken, RetusnSalesOrderController.FilterData)
	api.POST("", cekToken, RetusnSalesOrderController.Save)
	api.POST("/approve", cekToken, RetusnSalesOrderController.Approve)
	api.POST("/reject", cekToken, RetusnSalesOrderController.Reject)
	api.POST("/print/:id", cekToken, RetusnSalesOrderController.PrintReturnSO)
	api.POST("/payment/page/:page/count/:count", cekToken, RetusnSalesOrderController.FilterDataForSalesOrderReturn)

	api = r.Group("/api/return-sales-order-detail")
	api.POST("/page/:page/count/:count", ReturnSalesOrderDetailController.GetDetail)
	api.POST("", cekToken, ReturnSalesOrderDetailController.Save)
	api.POST("/updateQty", cekToken, ReturnSalesOrderDetailController.UpdateQty)
	api.DELETE("/:id", cekToken, ReturnSalesOrderDetailController.DeleteById)

	// RECEIVING
	api = r.Group("/api/receive")
	api.POST("/page/:page/count/:count", cekToken, ReceiveController.FilterData)
	api.GET("/:id", cekToken, ReceiveController.GetByReceiveId)
	api.POST("", cekToken, ReceiveController.Save)
	api.POST("/byPO", cekToken, ReceiveController.SaveByPO)
	api.POST("/reject", cekToken, ReceiveController.Reject)
	api.POST("/remove-PO", cekToken, ReceiveController.RemovePO)
	api.POST("/remove-PO-all", cekToken, ReceiveController.RemovePOAllItem)
	api.POST("/approve", cekToken, ReceiveController.Approve)
	api.POST("/print/:id", cekToken, ReceiveController.PrintPreview)
	api.POST("/export", cekToken, ReceiveController.Export)

	api = r.Group("/api/receive-detail")
	api.POST("/page/:page/count/:count", ReceiveDetailController.GetDetail)
	api.GET("/last-price/:productId", ReceiveDetailController.GetLastPrice)
	api.POST("", cekToken, ReceiveDetailController.Save)
	api.POST("/updateDetail", cekToken, ReceiveDetailController.UpdateDetail)
	api.DELETE("/:id", cekToken, ReceiveDetailController.DeleteByID)
	api.POST("/deleteMulti", cekToken, ReceiveDetailController.DeleteByIDMultiple)
	api.POST("/search-batch-expired/page/:page/count/:count", ReceiveDetailController.SearchBatchExpired)

	// RETURN RECEIVE
	api = r.Group("/api/return-receive")
	api.GET("/:id", cekToken, RetusnReceiveController.GetByReturnReceiveId)
	api.POST("/page/:page/count/:count", cekToken, RetusnReceiveController.FilterData)
	api.POST("", cekToken, RetusnReceiveController.Save)
	api.POST("/approve", cekToken, RetusnReceiveController.Approve)
	api.POST("/reject", cekToken, RetusnReceiveController.Reject)
	api.POST("/print/:id", cekToken, RetusnReceiveController.PrintReturn)

	// RETURN RECEIVE DETAIL
	api = r.Group("/api/return-receive-detail")
	api.POST("/page/:page/count/:count", cekToken, ReturnReceiveDetailController.GetDetail)
	api.POST("", cekToken, ReturnReceiveDetailController.Save)
	api.DELETE("/:id", cekToken, ReturnReceiveDetailController.DeleteById)

	// ADJUSTMENT
	api = r.Group("/api/adjustment")
	api.POST("/page/:page/count/:count", cekToken, AdjustmentController.FilterData)
	api.GET("/:id", cekToken, AdjustmentController.GetByAdjustmentId)
	api.POST("", cekToken, AdjustmentController.Save)
	api.POST("/approve", cekToken, AdjustmentController.Approve)
	api.POST("/print/:id", cekToken, AdjustmentController.PrintPreview)

	api = r.Group("/api/adjustment-detail")
	api.POST("/page/:page/count/:count", AdjustmentDetailController.GetDetail)
	api.POST("", cekToken, AdjustmentDetailController.Save)
	api.POST("/updateqty", cekToken, AdjustmentDetailController.UpdatQty)
	api.DELETE("/:id/:idAdj", cekToken, AdjustmentDetailController.DeleteByID)

	// PAYMENT
	api = r.Group("/api/payment")
	api.POST("/page/:page/count/:count", cekToken, PaymentController.FilterData)
	api.GET("/:id", cekToken, PaymentController.GetByPaymentId)
	api.POST("/salesOrderId/:salesOrderId", cekToken, PaymentController.GetBySalesOrderID)
	api.POST("", cekToken, PaymentController.Save)
	api.POST("/approve", cekToken, PaymentController.Approve)
	api.POST("/print/:id", cekToken, PaymentController.PrintPayment)

	api = r.Group("/api/payment-detail")
	api.POST("/all", cekToken, PaymentDetailController.GetDetail)
	api.POST("", cekToken, PaymentDetailController.Save)
	api.DELETE("/:id", cekToken, PaymentDetailController.DeleteById)

	api = r.Group("/api/payment-order")
	api.POST("/all", cekToken, PaymentOrderController.GetDetail)
	api.POST("", cekToken, PaymentOrderController.Save)
	api.DELETE("/:id", cekToken, PaymentOrderController.DeleteById)

	api = r.Group("/api/payment-return")
	api.POST("/all", cekToken, PaymentReturnController.GetDetail)
	api.POST("", cekToken, PaymentReturnController.Save)
	api.DELETE("/:id", cekToken, PaymentReturnController.DeleteById)

	// PAYMENT_DIRECT

	api = r.Group("/api/payment-direct")
	api.POST("/page/:page/count/:count", DirectPaymentController.FilterData)
	api.POST("/approve", DirectPaymentController.Approve)
	api.POST("/reject", DirectPaymentController.Reject)

	// Warehouse
	api = r.Group("/api/warehouse")
	api.GET("", WarehouseController.GetWarehouse)
	api.GET("/in", WarehouseController.GetWarehouseIn)
	api.GET("/out", WarehouseController.GetWarehouseOut)
	api.POST("/page/:page/count/:count", cekToken, WarehouseController.GetWarehouseFilter)
	api.GET("/id/:id", cekToken, WarehouseController.GetFilterWarehouse)
	api.POST("", cekToken, WarehouseController.SaveWarehouse)
	api.PUT("", cekToken, WarehouseController.UpdateWarehouse)
	api.GET("/like", cekToken, WarehouseController.GetWarehouseLike)

	api = r.Group("/api/salesman")
	api.GET("", cekToken, SalesmanController.GetSalesman)
	api.POST("/page/:page/count/:count", cekToken, SalesmanController.GetSalesmanFilter)
	api.GET("/id/:id", cekToken, SalesmanController.GetFilterSalesman)
	api.POST("", cekToken, SalesmanController.SaveSalesman)
	api.PUT("", cekToken, SalesmanController.UpdateSalesman)
	api.GET("/like", cekToken, SalesmanController.GetSalesmanLike)

	// STOCK - MUTATION
	api = r.Group("/api/stock-mutation")
	api.GET("/:id", StockMutationController.GetByStockMutationById)
	api.POST("/page/:page/count/:count", cekToken, StockMutationController.GetStockMutationPage)
	api.POST("", cekToken, StockMutationController.Save)
	api.POST("/approve", cekToken, StockMutationController.Approve)
	api.POST("/reject", cekToken, StockMutationController.Reject)
	api.POST("/print/:id", StockMutationController.PrintStockMutationForm)

	api = r.Group("/api/stock-mutation-detail")
	api.POST("/page/:page/count/:count", StockMutationDetailController.GetDetail)
	api.POST("", cekToken, StockMutationDetailController.Save)
	api.DELETE("/:id", cekToken, StockMutationDetailController.DeleteById)
	api.POST("/updateItemQty", cekToken, StockMutationDetailController.UpdateQty)

	// STOCK - MUTATION
	api = r.Group("/api/stock")
	api.GET("/:id/page/:page/count/:count", cekToken, StockController.GetByStockByProductIdPage)

	// HISTORY - STOCK
	api = r.Group("/api/history-stock")
	api.POST("/page/:page/count/:count", HistoryController.FilterDataHistoryStock)

	// PURCHASE ORDER
	api = r.Group("/api/purchase-order")
	api.POST("/page/:page/count/:count", cekToken, PurchaseOrderController.FilterData)
	api.GET("/:id", cekToken, PurchaseOrderController.GetByPurchaseOrderId)
	api.POST("", cekToken, PurchaseOrderController.Save)
	api.POST("/approve", cekToken, PurchaseOrderController.Approve)
	api.POST("/reject/:id", cekToken, PurchaseOrderController.Reject)
	api.POST("/cancel-submit/:id", cekToken, PurchaseOrderController.CancelSubmit)
	api.POST("/print/:id", cekToken, PurchaseOrderController.PrintPreview)
	api.POST("/print-by-pono/:pono", cekToken, PurchaseOrderController.PrintPreviewByPoNo)
	api.POST("/export", PurchaseOrderController.Export)

	api = r.Group("/api/purchase-order-detail")
	api.POST("/page/:page/count/:count", PurchaseOrderDetailController.GetDetail)
	api.POST("", cekToken, PurchaseOrderDetailController.Save)
	api.GET("/last-price/:productId", PurchaseOrderDetailController.GetLastPrice)
	api.DELETE("/:id", cekToken, PurchaseOrderDetailController.DeleteByID)
	api.POST("/updateDetail", cekToken, PurchaseOrderDetailController.UpdateDetail)

	// STOCK - OPNAME
	api = r.Group("/api/stock-opname")
	api.GET("/:id", StockOpnameController.GetByStockOpnameById)
	api.POST("/page/:page/count/:count", cekToken, StockOpnameController.GetStockOpnamePage)
	api.POST("", cekToken, StockOpnameController.Save)
	api.POST("/approve", cekToken, StockOpnameController.Approve)
	api.POST("/download-template/:warehouseId", StockOpnameController.DownloadTemplate)
	api.POST("/upload-template/:stockOpnameId", StockOpnameController.UploadTemplate)

	api = r.Group("/api/stock-opnames")
	api.GET("/recalulate-total/", StockOpnameController.RecalculateTotal)
	// api.POST("/reject", cekToken, StockOpnameController.Reject)
	// api.POST("/print/:id", StockOpnameController.PrintStockOpnameForm)

	api = r.Group("/api/stock-opname-detail")
	api.POST("/page/:page/count/:count", StockOpnameDetailController.GetDetail)
	api.POST("", cekToken, StockOpnameDetailController.Save)
	api.DELETE("/:id", cekToken, StockOpnameDetailController.DeleteById)
	api.POST("/updateItemQty", cekToken, StockOpnameDetailController.UpdateQty)

	// PAYMENT
	api = r.Group("/api/payment-supplier")
	api.POST("/page/:page/count/:count", cekToken, PaymentSupplierController.FilterData)
	api.GET("/:id", cekToken, PaymentSupplierController.GetByPaymentId)
	// api.POST("/salesOrderId/:salesOrderId", cekToken, PaymentController.GetBySalesOrderID)
	api.POST("", cekToken, PaymentSupplierController.Save)
	api.POST("/approve", cekToken, PaymentSupplierController.Approve)
	// api.POST("/print/:id", cekToken, PaymentController.PrintPayment)

	api = r.Group("/api/payment-supplier-detail")
	api.POST("/page/:page/count/:count", cekToken, PaymentSupplierDetailController.GetDetail)
	api.POST("", cekToken, PaymentSupplierDetailController.Save)
	api.DELETE("/:id", cekToken, PaymentSupplierDetailController.DeleteById)

	// REPORT
	api = r.Group("/api/report")
	api.POST("/payment-cash", cekToken, ReportPaymentCashController.DownloadReportPayment)
	api.POST("/payment-sales", cekToken, ReportPaymentCashController.DownloadReportSales)
	api.GET("/master-barang", cekToken, ReportMasterController.DownloadReportMasterProduct)
	api.POST("/payment-supplier", cekToken, ReportPaymentSupplierController.DownloadReporPaymentSupplier)

	// PARAMETER
	api = r.Group("/api/parameter")
	api.GET("/byname/:param-name", ParameterController.GetByName)
	api.GET("/all", ParameterController.GetAll)
	api.POST("", cekToken, ParameterController.UpdateParam)

	// -- END REPORT

	// Payment - SALES ORDER
	// api = r.Group("/api/payment")
	// api.GET("/:id", cekToken, RetusnSalesOrderController.GetByReturnSalesOrderId)
	// api.POST("/page/:page/count/:count", cekToken, RetusnSalesOrderController.FilterData)
	// api.POST("", cekToken, RetusnSalesOrderController.Save)
	// api.POST("/approve", cekToken, RetusnSalesOrderController.Approve)
	// api.POST("/reject", cekToken, RetusnSalesOrderController.Reject)
	// api.POST("/print/:id", cekToken, RetusnSalesOrderController.PrintReturnSO)

	// api = r.Group("/api/payment-detail")
	// api.POST("/page/:page/count/:count", ReturnSalesOrderDetailController.GetDetail)
	// api.POST("", cekToken, ReturnSalesOrderDetailController.Save)
	// api.POST("/updateQty", cekToken, ReturnSalesOrderDetailController.UpdateQty)
	// api.DELETE("/:id", cekToken, ReturnSalesOrderDetailController.DeleteById)

	// Dashboard
	dashboard := r.Group("/dashboard")
	dashboard.POST("/order-qty", DashboardController.FilterDataDashboard)

	apiVersion := r.Group("/api/version")
	apiVersion.GET("", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"date":    constants.VERSION_DATE,
			"version": constants.VERSION,
		})
	})

	return r

}

func cekSignature(c *gin.Context) {

	fmt.Println("cek signature")
	timestamp := c.Request.Header.Get("timestamp")
	signature := c.Request.Header.Get("signature")

	// 5H5GTtcehHqOLDgIzNu8
	key := beego.AppConfig.DefaultString("secret.key", "")
	// body := c.Request.Body
	body := "{}"
	res := dto.LoginResponseDto{}

	if ret := security.ValidateSignature(timestamp, key, signature, body); ret != true {
		res.ErrCode = kons.ERR_CODE_54
		res.ErrDesc = kons.ERR_CODE_54_MSG
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
	}
}

func cekToken(c *gin.Context) {

	res := models.Response{}
	tokenString := c.Request.Header.Get("Authorization")

	if strings.HasPrefix(tokenString, "Bearer ") == false {
		res.ErrCode = kons.ERR_CODE_54
		res.ErrDesc = kons.ERR_CODE_54_MSG
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			res.ErrCode = kons.ERR_CODE_54
			res.ErrDesc = kons.ERR_CODE_54_MSG
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(kons.TokenSecretKey), nil
	})

	if token != nil && err == nil {
		claims := token.Claims.(jwt.MapClaims)

		fmt.Println("claims : ", claims)

		fmt.Println("User name from TOKEN ", claims["user"])

		unixNano := time.Now().UnixNano()
		timeNowInInt := unixNano / 1000000

		tokenCreated := (claims["tokenCreated"])
		dto.CurrUser = (claims["user"]).(string)
		currUserId := (claims["userId"]).(string)
		dto.CurrUserId, _ = strconv.ParseInt(currUserId, 10, 64)

		fmt.Println("now : ", timeNowInInt)
		fmt.Println("token created time : ", tokenCreated)
		fmt.Println("user by token : ", dto.CurrUser)
		fmt.Println("user by token ID : ", dto.CurrUserId)

		tokenCreatedInString := tokenCreated.(string)
		tokenCreatedInInt, errTokenExpired := strconv.ParseInt(tokenCreatedInString, 10, 64)

		if errTokenExpired != nil {
			res.ErrCode = kons.ERR_CODE_55
			res.ErrDesc = kons.ERR_CODE_55_MSG
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		if ((timeNowInInt - tokenCreatedInInt) / 1000) > kons.TokenExpiredInMinutes {
			res.ErrCode = kons.ERR_CODE_55
			res.ErrDesc = kons.ERR_CODE_55_MSG
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}
		fmt.Println("Token already used for ", (timeNowInInt-tokenCreatedInInt)/1000, "sec, Max expired ", kons.TokenExpiredInMinutes, "sec ")
		// fmt.Println("token Valid ")

	} else {
		res.ErrCode = kons.ERR_CODE_54
		res.ErrDesc = kons.ERR_CODE_54_MSG
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}
}

// CORSMiddleware ...
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if c.Request.Method == "OPTIONS" {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
// 			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
// 			c.Writer.Header().Set("Content-Type", "application/json, charset=utf-8")
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

// // @APIVersion 1.0.0
// // @Title beego Test API
// // @Description beego has a very cool tools to autogenerate documents for your API
// // @Contact astaxie@gmail.com
// // @TermsOfServiceUrl http://beego.me/
// // @License Apache 2.0
// // @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
// package routers

// import (
// 	"distribution-system-be/controllers"

// 	"github.com/astaxie/beego"
// )

// func init() {
// 	ns := beego.NewNamespace("/v1",

// 		beego.NSNamespace("/object",
// 			beego.NSInclude(
// 				&controllers.ObjectController{},
// 			),
// 		),
// 		beego.NSNamespace("/user",
// 			beego.NSInclude(
// 				&controllers.UserController{},
// 			),
// 		),
// 	)
// 	beego.AddNamespace(ns)
// }
