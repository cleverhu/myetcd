package services

import (
	"context"
	"fmt"
	"net/http"
)

func ProdEncodeFunc(ctx context.Context, request *http.Request, requestParam interface{}) error {
	prod := requestParam.(ProdRequest)
	request.URL.Path += "/prod/" + fmt.Sprintf("%d", prod.Id)
	return nil
}
