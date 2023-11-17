package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/performance"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
)

type HeadlessBrowser struct {
	url     string
	headers map[string]interface{}
}

func NewHeadlessBrowser(url string, headers []string, cookie string) HeadlessBrowser {
	h := make(map[string]interface{})
	for _, header := range headers {
		parts := strings.Split(header, ":")
		if len(parts) == 2 {
			h[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	h["Cookie"] = cookie

	return HeadlessBrowser{
		url:     url,
		headers: h,
	}
}

// Screenshot 指定したURLのスクリーンショットを取得する
func (b HeadlessBrowser) Screenshot(waitSec int) (*[]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if e := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(b.url),
		network.SetExtraHTTPHeaders(b.headers),
		chromedp.Sleep(time.Duration(waitSec) * time.Second),
		chromedp.CaptureScreenshot(&buf),
	}); e != nil {
		return nil, e
	}

	return &buf, nil
}

type CoreWebVital struct {
	LargestContentfulPaint float64
	FirstInputDelay        float64
	CumulativeLayoutShift  float64
	TimeToFirstByte        float64
	FirstContentfulPaint   float64
	TotalBlockingTime      float64
	Interactive            float64
}

// GetCoreWebVital CoreWebVitalを計測する
// HACK: この機能ではCoreWebVitalは取得できないので変更が必要
func (b HeadlessBrowser) GetCoreWebVital(waitSec int) (CoreWebVital, error) {
	type Entry struct {
		Name          string  `json:"name"`
		EntryType     string  `json:"entryType"`
		StartTime     float64 `json:"startTime"`
		Duration      float64 `json:"duration"`
		InitiatorType string  `json:"initiatorType"`
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var vital CoreWebVital
	var r float64

	if e := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(b.url),
		chromedp.Sleep(2 * time.Second), // ページが読み込まれるのを待つ
		chromedp.ActionFunc(func(ctx context.Context) error {
			var metrics string
			if e := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`JSON.stringify(performance.getEntries())`, &metrics)); e != nil {
				return e
			}

			var entries []Entry
			if e := json.Unmarshal([]byte(metrics), &entries); e != nil {
				return e
			}

			return nil
		}),
	}); e != nil {
		return vital, e
	}

	fmt.Println(r)

	return vital, nil
}

/*
getMetric(ctx, "largest_contentful_paint", &vital.LargestContentfulPaint),
		getMetric(ctx, "first_input_delay", &vital.FirstInputDelay),
		getMetric(ctx, "cumulative_layout_shift", &vital.CumulativeLayoutShift),
		getMetric(ctx, "time_to_first_byte", &vital.TimeToFirstByte),
		getMetric(ctx, "first_contentful_paint", &vital.FirstContentfulPaint),
		getMetric(ctx, "total_blocking_time", &vital.TotalBlockingTime),
		getMetric(ctx, "interactive", &vital.Interactive),
*/

// getMetric 指定したCoreWebVitalの指標を取得する
// ちょっと効率悪い気がする... まぁCLIなのでいいかな...
func getMetric(ctx context.Context, metricName string, value *float64) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		metrics, err := performance.GetMetrics().Do(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(metrics)
		for _, metric := range metrics {
			fmt.Println(metric.Name)
			if metric.Name == metricName {
				*value = metric.Value
				break
			}
		}
		return nil
	})
}
