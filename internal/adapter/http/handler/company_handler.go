package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	company "main.go/internal/usecase/companies"
)

type CompanyHandler struct {
	getByIDUC   *company.GetCompanyByIDUseCase
	getNearbyUC *company.GetNearbyCompaniesUseCase
	searchUC    *company.SearchCompaniesUseCase
}

func NewCompanyHandler(
	getByIDUC *company.GetCompanyByIDUseCase,
	getNearbyUC *company.GetNearbyCompaniesUseCase,
	searchUC *company.SearchCompaniesUseCase,
) *CompanyHandler {
	return &CompanyHandler{
		getByIDUC:   getByIDUC,
		getNearbyUC: getNearbyUC,
		searchUC:    searchUC,
	}
}

// GET /v1/companies
// Example: ?lat=-23.5&lng=-46.6&radius=10
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")
	keyword := c.Query("keyword")
	category := c.Query("category")
	fetchAllStr := c.DefaultQuery("fetch_all", "false")
	pageToken := c.Query("page_token")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lng, errLng := strconv.ParseFloat(lngStr, 64)
	radius, _ := strconv.ParseFloat(radiusStr, 64)
	fetchAll, _ := strconv.ParseBool(fetchAllStr)

	if errLat != nil || errLng != nil || lat == 0 || lng == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "latitude and longitude are required and must be valid numbers",
		})

		return
	}

	if radius <= 0 {
		radius = 800.0
	}

	if keyword != "" || category != "" {
		input := company.SearchCompaniesInput{
			Latitude:  lat,
			Longitude: lng,
			RadiusKm:  radius,
			Keyword:   keyword,
			Category:  category,
		}

		companies, err := h.searchUC.Execute(c.Request.Context(), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  companies,
			"total": len(companies),
		})
		return
	}

	input := company.GetNearbyInput{
		Latitude:  lat,
		Longitude: lng,
		RadiusKm:  radius,
		FetchAll:  fetchAll,
		PageToken: pageToken,
	}

	companies, err := h.getNearbyUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":  companies,
		"total": len(companies),
	}

	c.JSON(http.StatusOK, response)
}

// GET /v1/companies/:id
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	companyEntity, err := h.getByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": companyEntity,
	})
}
