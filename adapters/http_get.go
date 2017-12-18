package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/smartcontractkit/chainlink-go/models"
)

type HttpGet struct {
	AdapterBase
	Endpoint string `json:"endpoint"`
}

func (self *HttpGet) Perform(input models.RunResult) models.RunResult {
	response, err := http.Get(self.Endpoint)
	if err != nil {
		return models.RunResult{Error: err}
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return models.RunResult{Error: err}
	}
	if response.StatusCode >= 300 {
		return models.RunResult{Error: fmt.Errorf(body)}
	}

	return models.RunResultWithValue(body)
}
