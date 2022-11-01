package dot

import (
	"context"

	"github.com/weblazy/easy/utils/elog"
	"github.com/weblazy/easy/utils/http/http_client"
	"github.com/weblazy/easy/utils/http/http_client/http_client_config"
	"go.uber.org/zap"
)

// DotRequest request polkadot
func DotRequest(ctx context.Context, req interface{}) ([]byte, error) {
	cfg := http_client_config.DefaultConfig()
	client := http_client.NewHttpClient(cfg)
	request := client.Request.SetContext(ctx)

	resp, err := request.
		ExpectContentType("application/json").
		SetBody(req).
		Post("https://rpc.polkadot.io")
	if err != nil {
		elog.ErrorCtx(ctx, "DotRequestErr", zap.Error(err))
		return nil, err
	}
	return resp.Body(), nil
}
