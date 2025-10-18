package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkiar00/homework/internal/usecase"
)

// EchoController translates HTTP requests to usecase calls.
type EchoController struct{
	uc *usecase.EchoUsecase
}

func NewEchoController(uc *usecase.EchoUsecase) *EchoController {
	return &EchoController{uc: uc}
}

// GetEcho handles GET /echo/:value
func (ec *EchoController) GetEcho(c *gin.Context) {
	v := c.Param("value")
	res, err := ec.uc.Echo(v)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": res})
}

// PostEcho handles POST /echo with JSON body {"value":"..."}
func (ec *EchoController) PostEcho(c *gin.Context) {
	var body struct{ Value string `json:"value"` }
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	res, err := ec.uc.Echo(body.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": res})
}
