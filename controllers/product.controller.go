package controllers

import (
	"distribution-system-be/constants"
	"distribution-system-be/database"
	"distribution-system-be/models"
	dbmodels "distribution-system-be/models/dbModels"
	dto "distribution-system-be/models/dto"
	"distribution-system-be/services"
	"distribution-system-be/utils/util"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ProductController ...
type ProductController struct {
	DB *gorm.DB
}

//ProductService ...
var ProductService = new(services.ProductService)

//GetProductListPaging ...
func (h *ProductController) GetProductListPaging(c *gin.Context) {
	req := dto.FilterProduct{}
	res := models.ResponsePagination{}

	page, errPage := strconv.Atoi(c.Param("page"))
	if errPage != nil {
		logs.Info("error", errPage)
		res.Error = errPage.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	count, errCount := strconv.Atoi(c.Param("count"))
	if errCount != nil {
		logs.Info("error", errPage)
		res.Error = errCount.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ProductService.GetProductFilterPaging(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

//GetProductListPaging ...
func (h *ProductController) GetProductListPagingAllStatus(c *gin.Context) {
	req := dto.FilterProduct{}
	res := models.ResponsePagination{}

	page, errPage := strconv.Atoi(c.Param("page"))
	if errPage != nil {
		logs.Info("error", errPage)
		res.Error = errPage.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	count, errCount := strconv.Atoi(c.Param("count"))
	if errCount != nil {
		logs.Info("error", errPage)
		res.Error = errCount.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ProductService.GetProductFilterPagingAllStatus(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

func (h *ProductController) SearchProduct(c *gin.Context) {
	req := dto.FilterProduct{}
	res := models.ResponsePagination{}

	page, errPage := strconv.Atoi(c.Param("page"))
	if errPage != nil {
		logs.Info("error", errPage)
		res.Error = errPage.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	count, errCount := strconv.Atoi(c.Param("count"))
	if errCount != nil {
		logs.Info("error", errPage)
		res.Error = errCount.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.Error = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ProductService.SearchProduct(req, page, count)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// GetProductDetails ...
func (h *ProductController) GetProductDetails(c *gin.Context) {
	res := models.ContentResponse{}

	id, errID := strconv.Atoi(c.Param("id"))
	if errID != nil {
		logs.Info("error", errID)
		// res.Error = errID.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	res = ProductService.GetProductDetails(id)

	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

// SaveProduct ...
func (h *ProductController) SaveProduct(c *gin.Context) {

	req := dbmodels.Product{}
	res := models.NoContentResponse{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = "03"
		res.ErrDesc = "Error, unmarshall body Request"
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, ProductService.SaveProduct(&req))
	return
}

// UpdateProduct ...
func (h *ProductController) UpdateProduct(c *gin.Context) {
	req := dbmodels.Product{}
	res := models.NoContentResponse{}

	body := c.Request.Body
	dataBodyReq, _ := ioutil.ReadAll(body)

	if err := json.Unmarshal(dataBodyReq, &req); err != nil {
		fmt.Println("Error, body Request ")
		res.ErrCode = "03"
		res.ErrDesc = "Error, unmarshall body Request"
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, ProductService.UpdateProduct(&req))
	return
}

func (h *ProductController) ProductList(c *gin.Context) {
	// res := []dbmodels.Product{}

	c.JSON(http.StatusOK, ProductService.ProductList())
	return
}

//GetProductLike ...
func (h *ProductController) GetProductLike(c *gin.Context) {
	res := models.ContentResponse{}

	productterms := c.Query("terms")

	if productterms == "" {
		logs.Info("error", "can't found the brand string")
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	// fmt.Sprintf("ini lookupstr = " + lookupstr)

	res = ProductService.GetProductLike(productterms)
	c.JSON(http.StatusOK, res)
	c.Abort()
	return
}

func (h *ProductController) ProcessCSV(c *gin.Context) {

	fmt.Println("Process")
	fileObat, err := os.Open("obat.csv")
	if err != nil {
		fmt.Println("Error ==>", err)
		panic(err)
	}
	defer fileObat.Close()

	fmt.Println("reader")
	csvReader := csv.NewReader(fileObat)
	csvReader.Comma = ';'
	records, err2 := csvReader.ReadAll()
	if err2 != nil {
		fmt.Println("Error err2 ==>", err2)
		return
	}

	var datas []dto.TemplateObat
	for _, value := range records {
		fmt.Println("data ", value)
		var template dto.TemplateObat
		template.Plu = value[0]
		template.Name = value[1]
		template.Satuan = value[2]
		template.Qty = util.Atoi64(value[3])
		template.Hargabeli = math.Floor(util.AtoFloat64(value[7]) / util.AtoFloat64(value[3]))
		template.Hargajual = util.Atoi64(value[6])
		template.SatuanBesar = value[4]
		template.Kadar = value[8]
		datas = append(datas, template)

	}

	for _, templateObat := range datas {
		fmt.Println(templateObat)
		var produk dbmodels.Product
		var uomID int64
		lookup, errcode, _, _ := database.GetLookupByName(templateObat.Satuan)
		if errcode != constants.ERR_CODE_00 {
			uomID = 30
		} else {
			uomID = lookup.ID
		}

		var uomID2 int64
		lookup2, errcode2, _, _ := database.GetLookupByName(templateObat.SatuanBesar)
		if errcode2 != constants.ERR_CODE_00 {
			uomID2 = 30
		} else {
			uomID2 = lookup2.ID
		}

		produk.BigUomID = uomID
		produk.BrandID = 1
		produk.Code = ""
		produk.Hpp = float32(templateObat.Hargabeli)
		produk.LastUpdate = time.Now()
		produk.LastUpdateBy = "system"
		produk.Name = templateObat.Name
		produk.ProductGroupID = 1
		produk.PLU = templateObat.Plu
		produk.QtyUom = int16(templateObat.Qty)
		produk.SellPrice = float32(templateObat.Hargajual)
		produk.SellPriceType = 0
		produk.SmallUomID = uomID
		produk.BigUomID = uomID2
		produk.Status = 1
		produk.Composition = templateObat.Kadar

		database.SaveProduct(produk)
	}

	c.JSON(http.StatusOK, "ok")
	c.Abort()
	return
}

func (h *ProductController) ProcessUpdateProd(c *gin.Context) {

	fmt.Println("Process")
	fileObat, err := os.Open("update-product.csv")
	if err != nil {
		fmt.Println("Error ==>", err)
		panic(err)
	}
	defer fileObat.Close()

	fmt.Println("reader")
	csvReader := csv.NewReader(fileObat)
	csvReader.Comma = ';'
	records, err2 := csvReader.ReadAll()
	if err2 != nil {
		fmt.Println("Error err2 ==>", err2)
		return
	}

	var datas []dto.TemplateObat
	i := 0
	for _, value := range records {
		fmt.Println("data ", value)
		var template dto.TemplateObat
		template.Plu = value[0]
		template.Name = value[1]
		template.Satuan = value[2]
		template.HargaJualBaru = util.AtoFloat32(value[3])

		datas = append(datas, template)
		i++
		if i > 10 {
			continue
		}
	}

	for _, templateObat := range datas {
		fmt.Println(templateObat)
		// var produk dbmodels.Product
		if templateObat.Plu == "" {
			continue
		}

		produk, errCode, _ := database.FindProductByPLU(templateObat.Plu)
		if errCode != constants.ERR_CODE_00 {
			continue
		}
		fmt.Println("Product ", produk)
		database.UpdateProductByPLU(produk.ID, templateObat.HargaJualBaru)
	}

	c.JSON(http.StatusOK, "ok")
	c.Abort()
	return
}
