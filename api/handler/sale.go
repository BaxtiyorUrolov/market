package handler

import (
	"context"
	"fmt"
	"market/api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateSale godoc
// @Router       /sale [POST]
// @Summary      Create a new sale
// @Description  create a new sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 sale body models.CreateSale false "sale"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateSale(c *gin.Context) {
	sale := models.CreateSale{}
	if err := c.ShouldBindJSON(&sale); err != nil {
		handleResponse(c, h.log, "error is while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Sale().Create(context.Background(), sale)
	if err != nil {
		handleResponse(c, h.log, "error is while creating sale", http.StatusInternalServerError, err.Error())
		return
	}

	createdBranch, err := h.storage.Sale().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, h.log, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusCreated, createdBranch)
}

// GetSale godoc
// @Router       /sale/{id} [GET]
// @Summary      Get sale by id
// @Description  get sale by id
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 id path string true "sale_id"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetSale(c *gin.Context) {
	uid := c.Param("id")

	sale, err := h.storage.Sale().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, h.log, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, sale)
}

// GetSaleList godoc
// @Router       /sales [GET]
// @Summary      Get sale list
// @Description  get sale list
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetSaleList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error is while converting page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	sales, err := h.storage.Sale().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, h.log, "Error is while getting Sale list: ", http.StatusInternalServerError, err.Error())
		return 
	}

	handleResponse(c, h.log, "", http.StatusOK, sales)
}

// UpdateSale godoc
// @Router       /sale/{id} [PUT]
// @Summary      Update sale
// @Description  Update sale by ID
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param        id path string true "Sale ID"
// @Param        sale body models.UpdateSale true "Sale object"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
func (h Handler) UpdateSale(c *gin.Context) {
    uid := c.Param("id")
    sale := models.UpdateSale{}

    if uid == "" {
        handleResponse(c, h.log, "error is while reading body", http.StatusBadRequest, "sale ID is empty")
        return
    }

    if err := c.ShouldBindJSON(&sale); err != nil {
        handleResponse(c, h.log, "error is while reading body", http.StatusBadRequest, err.Error())
        return
    }

	status, err := h.storage.Sale().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, h.log, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	if status.Status != "in_process" {
		handleResponse(c, h.log, "sale status is not 'in_process', cannot update", http.StatusBadRequest, "")
		return
	}

    sale.ID = uid

    baskets, err := h.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page: 1,
		Limit: 100,  
		Search: uid,
    })
    if err != nil {
        handleResponse(c, h.log, "error is while getting baskets by sale ID", http.StatusInternalServerError, err.Error())
        return
    }

    totalPrice := 0
    for _, basket := range baskets.Baskets {
        totalPrice += basket.Price
    }

	fmt.Println(totalPrice)

	count, err := h.storage.Repository().ProductByID(context.Background(), baskets.Baskets[0].ProductID)
	if err != nil {
		return 
	}

	if count - baskets.Baskets[0].Quantity < 0 {
		return
	}

    sale.Price = float32(totalPrice)

    id, err := h.storage.Sale().Update(context.Background(), sale)
    if err != nil {
        handleResponse(c, h.log, "error is while updating sale", http.StatusInternalServerError, err.Error())
        return
    }

    updatedSale, err := h.storage.Sale().GetByID(context.Background(), id)
    if err != nil {
        handleResponse(c, h.log, "error is while getting by ID", http.StatusInternalServerError, err.Error())
        return
    }

    handleResponse(c, h.log, "", http.StatusOK, updatedSale)
}

// DeleteSale godoc
// @Router       /sale/{id} [DELETE]
// @Summary      Delete sale
// @Description  delete sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 id path string true "sale_id"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteSale(c *gin.Context) {
	uid := c.Param("id")
	if err := h.storage.Sale().Delete(context.Background(), uid); err != nil {
		handleResponse(c, h.log, "error is while deleting", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "sale deleted!")
}
