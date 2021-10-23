package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// CSAKeysController manages CSA keys
type CSAKeysController struct {
	controller
}

func NewCSAKeysController(app chainlink.Application) CSAKeysController {
	return CSAKeysController{newController(app)}
}

// Index lists CSA keys
// Example:
// "GET <application>/keys/csa"
func (ctrl *CSAKeysController) Index(c *gin.Context) {
	keys, err := ctrl.app.GetKeyStore().CSA().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewCSAKeyResources(keys), "csaKeys")
}

// Create and return a CSA key
// Example:
// "POST <application>/keys/csa"
func (ctrl *CSAKeysController) Create(c *gin.Context) {
	key, err := ctrl.app.GetKeyStore().CSA().Create()
	if err != nil {
		if errors.Is(err, keystore.ErrCSAKeyExists) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewCSAKeyResource(key), "csaKeys")
}

// Exports a key
func (ctrl *CSAKeysController) Export(c *gin.Context) {
	defer ctrl.lggr.ErrorIfClosing(c.Request.Body, "Export request body")

	keyID := c.Param("keyID")
	newPassword := c.Query("newpassword")

	bytes, err := ctrl.app.GetKeyStore().CSA().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, MediaType, bytes)
}
