package httputil

import (
	"context"
	"time"

	"github.com/huaxr/framework/pkg/confutil"
	"github.com/spf13/cast"

	"github.com/huaxr/framework/internal/define"
	"github.com/huaxr/framework/logx"

	"github.com/go-resty/resty/v2"
)

type HttpCli struct {
	client  *resty.Client
	timeout time.Duration
	retry   int
	// 是否需要透传字段
	needPt bool
}

func NewHttpClient(timeout time.Duration, retry int, needPt bool) *HttpCli {
	r := resty.New()
	if timeout > 0 {
		r.SetTimeout(timeout)
	}

	if retry > 0 {
		r.SetRetryCount(retry)
	}

	r.RetryHooks = []resty.OnRetryFunc{
		func(r *resty.Response, err error) {
			logx.T(nil, define.ArchError).Infof("Retry happens url:%v, err:%v",
				r.Header().Get("host"), err)
		},
	}
	return &HttpCli{
		client:  r,
		timeout: timeout,
		retry:   retry,
		needPt:  needPt,
	}
}

// 透传信息  gin->gin (user->gin from middleware)
func (h *HttpCli) RequestWithCtx(ctx context.Context) *resty.Request {
	var req = h.client.R().SetContext(ctx)
	for _, i := range define.CtxPassThrough {
		switch i {
		case define.CallFrom:
			req = req.SetHeader(i.String(), confutil.GetDefaultConfig().PSM.String())
		default:
			if v := ctx.Value(i.String()); v != nil {
				req = req.SetHeader(i.String(), cast.ToString(v))
			}
		}
	}
	return req
}

// if you wanna use form-data post, RequestWithCtx is helpful by setting your own headers.
func (h *HttpCli) Post(ctx context.Context, url string, data interface{}) ([]byte, error) {
	// token: SetHeader("Authorization", fmt.Sprintf("Bearer %s", req.Token)).
	res, err := h.RequestWithCtx(ctx).SetHeader("Content-Type", "application/json").SetBody(data).Post(url)
	if err != nil {
		logx.T(ctx, define.ArchError).Errorf("post:%s err:%v", url, err)
		return nil, err
	}
	// cause err catch already, <-ctx.Done() dose not reached in select
	return res.Body(), nil
}

func (h *HttpCli) Get(ctx context.Context, url string) ([]byte, error) {
	res, err := h.RequestWithCtx(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	return res.Body(), nil
}
