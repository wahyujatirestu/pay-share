package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/dto"
	"github.com/wahyujatirestu/payshare/model"
	"github.com/wahyujatirestu/payshare/service"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (c *ProductController) Create(ctx *gin.Context) {
	var req dto.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	product := &model.Product{
		Name: 			req.Name,
		Description: 	req.Description,
		Price: 			req.Price,
		Unit: 			req.Unit,
	}

	if err := c.productService.Create(product); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Product created successfully", "product": product,
	})
}

func (c *ProductController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := c.productService.GetById(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if product == nil {
		ctx.JSON(404, gin.H{"error": "Product not found"})
		return 
	}

	ctx.JSON(200, gin.H{"product": product})
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	filters := make(map[string]interface{})

	if name := ctx.Query("name"); name != "" {
		filters["name"] = name
	}

	if priceMin := ctx.Query("price_min"); priceMin != "" {
		filters["price_min"] = priceMin
	}

	if priceMax := ctx.Query("price_max"); priceMax != "" {
		filters["price_max"] = priceMax
	}

	if unit := ctx.Query("unit"); unit != "" {
		filters["unit"] = unit
	}

	products, err := c.productService.GetAll(filters)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"products": products})
}

func (c *ProductController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var req dto.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	product := &model.Product{
		ID:        		id,
		Name:      		req.Name,
		Description: 	req.Description,	
		Price:     		req.Price,
		Unit:      		req.Unit,
	}

	if err := c.productService.Update(product); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
	}
	ctx.JSON(200, gin.H{"message": "Product updated successfully", "product": product})
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.productService.Delete(id); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Product deleted successfully"})
}