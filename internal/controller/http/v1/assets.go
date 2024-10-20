package v1

import (
	"bhs/internal/entity"
	"bhs/internal/usecase"
	"bhs/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type assetsRoutes struct {
	a  usecase.Assets
	au usecase.Auth
	l  logger.Interface
}

type assetsResponse struct {
	Assets []entity.Asset `json:"assets"`
}

func newAssetsRoutes(handler *gin.RouterGroup, a usecase.Assets, au usecase.Auth, l logger.Interface) {
	r := &assetsRoutes{a, au, l}

	h := handler.Group("/assets")
	{
		h.GET("", r.listAssets)
		h.GET("/lib", r.listUserAssets)
		h.GET("/buy", r.buyAsset)
	}
}

func (r *assetsRoutes) listAssets(c *gin.Context) {
	pageNumQ, foundPageNum := c.GetQuery("page")
	pageItemsCountQ, foundPageItemsCount := c.GetQuery("per_page")
	println(pageNumQ)
	println(pageItemsCountQ)
	if !foundPageNum {
		pageNumQ = "1"
	}
	if !foundPageItemsCount {
		pageItemsCountQ = "10"
	}

	var assets []entity.Asset
	page, _ := strconv.ParseUint(pageNumQ, 15, 64)
	pageItemsCount, _ := strconv.ParseUint(pageItemsCountQ, 15, 64)
	println(page)
	println(pageItemsCount)
	assets, err := r.a.GetAssetsPage(c, page, pageItemsCount)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, assetsResponse{assets})
}

func (r *assetsRoutes) listUserAssets(c *gin.Context) {
	pageNumQ, foundPageNum := c.GetQuery("page")
	pageItemsCountQ, foundPageItemsCount := c.GetQuery("per_page")

	rawAuthToken := c.GetHeader("Authorization")
	token, validFormat := strings.CutPrefix(rawAuthToken, "Bearer ")
	if !validFormat {
		errorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
		return
	}
	user, err := r.au.Authorize(c, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "invalid authorization token")
		return
	}
	var assets []entity.Asset
	if foundPageNum && foundPageItemsCount {
		page, _ := strconv.ParseUint(pageNumQ, 15, 64)
		pageItemsCount, _ := strconv.ParseUint(pageItemsCountQ, 15, 64)
		assets, err = r.a.GetUserAssetsPage(c, user, page, pageItemsCount)
	} else {
		assets, err = r.a.GetUserAssets(c, user)
	}
	if err != nil {
		r.l.Error(err)
		errorResponse(c, http.StatusBadRequest, "cannot load assets")
	} else {
		c.JSON(http.StatusOK, assetsResponse{assets})
	}
}

func (r *assetsRoutes) buyAsset(c *gin.Context) {
	rawAuthToken := c.GetHeader("Authorization")
	token, validFormat := strings.CutPrefix(rawAuthToken, "Bearer ")
	if !validFormat {
		errorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
		return
	}
	user, err := r.au.Authorize(c, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "invalid authorization token")
		return
	}
	assetId, foundQ := c.GetQuery("asset_id")
	assetIdInt, err := strconv.ParseInt(assetId, 0, 15)
	if !foundQ || assetIdInt == 0 || err != nil {
		errorResponse(c, http.StatusUnauthorized, "'asset_id' query param not provided")
	}
	ok, err := r.a.AddUserAsset(c, user, assetIdInt)
	if !ok || err != nil {
		errorResponse(c, http.StatusUnauthorized, "asset cannot be added")
	}
	c.Status(http.StatusOK)
}
