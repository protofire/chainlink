package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// VRFKeysController manages VRF keys
type VRFKeysController struct {
	controller
}

func NewVRFKeysController(app chainlink.Application) VRFKeysController {
	return VRFKeysController{newController(app)}
}

// Index lists VRF keys
// Example:
// "GET <application>/keys/vrf"
func (vrfkc *VRFKeysController) Index(c *gin.Context) {
	keys, err := vrfkc.app.GetKeyStore().VRF().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResources(keys), "vrfKey")
}

// Create and return a VRF key
// Example:
// "POST <application>/keys/vrf"
func (vrfkc *VRFKeysController) Create(c *gin.Context) {
	pk, err := vrfkc.app.GetKeyStore().VRF().Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResource(pk), "vrfKey")
}

// Delete a VRF key
// Example:
// "DELETE <application>/keys/vrf/:keyID"
// "DELETE <application>/keys/vrf/:keyID?hard=true"
func (vrfkc *VRFKeysController) Delete(c *gin.Context) {
	keyID := c.Param("keyID")
	key, err := vrfkc.app.GetKeyStore().VRF().Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = vrfkc.app.GetKeyStore().VRF().Delete(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResource(key), "vrfKey")
}

// Import imports a VRF key
// Example:
// "Post <application>/keys/vrf/import"
func (vrfkc *VRFKeysController) Import(c *gin.Context) {
	defer vrfkc.lggr.ErrorIfClosing(c.Request.Body, "Import request body")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := vrfkc.app.GetKeyStore().VRF().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewVRFKeyResource(key), "vrfKey")
}

// Export exports a VRF key
// Example:
// "Post <application>/keys/vrf/export/:keyID"
func (vrfkc *VRFKeysController) Export(c *gin.Context) {
	defer vrfkc.lggr.ErrorIfClosing(c.Request.Body, "Export request body")

	keyID := c.Param("keyID")
	// New password to re-encrypt the export with
	newPassword := c.Query("newpassword")
	bytes, err := vrfkc.app.GetKeyStore().VRF().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
