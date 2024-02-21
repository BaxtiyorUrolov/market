package handler

import (
	"context"
	"net/http"
	"strconv"
	"market/api/models"

	"github.com/gin-gonic/gin"
)

// CreateBasket godoc
// @Router       /basket [POST]
// @Summary      Create a new basket
// @Description  create a new basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 basket body models.CreateBasket false "basket"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *Handler) CreateBasket(c *gin.Context)  {
	basket := models.CreateBasket{}
	

	if err :=  c.ShouldBindJSON(&basket); err != nil {
		handleResponse(c, h.log, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	resp, err :=  h.services.Basket().Create(context.Background(), basket)
	if err != nil {
		handleResponse(c, h.log, "error is while creating basket", http.StatusInternalServerError, err.Error())
		return 
	}

	handleResponse(c, h.log, "", http.StatusCreated, resp)
}

// GetBasket godoc
// @Router       /basket/{id} [GET]
// @Summary      Get basket by id
// @Description  get basket by id
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 id path string true "basket_id"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBasket(c *gin.Context) {
	uid := c.Param("id")

	basket, err := h.services.Basket().Get(context.Background(), uid)
	if err != nil {
		handleResponse(c, h.log, "error is while getting by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, basket)
}

// GetBasketList godoc
// @Router       /baskets [GET]
// @Summary      Get basket list
// @Description  get basket list
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.BasketsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBasketList(c *gin.Context) {
	var (
		page, limit int
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, h.log, "error while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, h.log, "error while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search := c.Query("search")

	response, err := h.services.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, h.log, "error while getting basket list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, response)
}

// UpdateBasket godoc
// @Router       /basket/{id} [PUT]
// @Summary      Update basket
// @Description  get basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 id path string true "basket_id"
// @Param 		 basket body models.UpdateBasket false "basket"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateBasket(c *gin.Context) {
	updatedBasket := models.UpdateBasket{}

	uid := c.Param("id")
	if err := c.ShouldBindJSON(&updatedBasket); err != nil {
		handleResponse(c, h.log, "error is while decoding ", http.StatusBadRequest, err)
		return
	}

	updatedBasket.ID = uid

	basket, err := h.services.Basket().Update(context.Background(), updatedBasket)
	if err != nil {
		handleResponse(c, h.log, "error is while updating basket", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, basket)
}

// DeleteBasket godoc
// @Router       /basket/{id} [DELETE]
// @Summary      Delete basket
// @Description  delete basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 id path string true "basket_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteBasket(c *gin.Context) {
	uid := c.Param("id")

	if err := h.services.Basket().Delete(context.Background(), models.PrimaryKey{ID: uid}); err != nil {
		handleResponse(c, h.log, "error is while deleting basket", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, h.log, "", http.StatusOK, nil)
}