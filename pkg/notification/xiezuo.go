package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	gotemplate "html/template"
	"strings"
	"time"

	"github.com/prometheus/alertmanager/template"
	"github.com/shaowenchen/grafana-webhook/config"
	"github.com/shaowenchen/grafana-webhook/pkg/http"
)

func SendXieZuo(body []byte) {

	nfData := map[string]interface{}{
		"content": GenerateBody(body),
	}
	req := map[string]interface{}{
		"msgtype": "text",
		"text":    nfData,
	}

	_, err := http.Post(req, config.Config.Notify.XieZuo)
	if err != nil {
		fmt.Printf("[SendXieZuo]push server error,[%v]\n", err)
	}
}

func GenerateBody(body []byte) string {
	var grafanaData ExtendedData
	err := json.Unmarshal(body, &grafanaData)
	if err != nil {
		fmt.Printf("[GenerateBody]unmarshal error,[%v]\n", err)
	}
	tplData := TemplateData{
		AlertLabels: " 无",
		AlertValues: " NoData",
	}
	for _, item := range grafanaData.Alerts {
		if item.Status == "firing" {
			test, _ := json.Marshal(item)
			fmt.Println(string(test))
			tplData.AlertURL = ""
			var cstZone = time.FixedZone("CST", 8*3600) 
			tplData.AlertAtTime = item.StartsAt.In(cstZone).Format("2006-01-02 15:04:05")
			if item.Labels["rulename"] != "" {
				tplData.AlertName = item.Labels["rulename"]
			} else {
				tplData.AlertName = item.Labels["alertname"]
				metricsItem := parseValueString(item.ValueString)
				tplData.AlertLabels = "\n " + strings.Join(metricsItem[0].Labels, "\n")
				tplData.AlertValues = metricsItem[0].Value
			}
			tplData.AlertKeepTime = time.Now().Sub(item.StartsAt).String()
			goto A
		}
	}
A:

	var buf bytes.Buffer
	tmpl := gotemplate.New("xiezuo")
	tmpl.Parse(config.Config.Notify.Template)
	err = tmpl.Execute(&buf, tplData)
	if err != nil {
		fmt.Printf("[GenerateBody]template error,[%v]\n", err)
	}
	fmt.Println(buf.String())
	return buf.String()
}

type ValuesStringItem struct {
	Metric string   `json:"metric"`
	Value  string   `json:"value"`
	Labels []string `json:"labels"`
}

func parseValueString(valueString string) []ValuesStringItem {
	metrics_str := strings.Split(valueString, "], [")
	resultsValues := make([]ValuesStringItem, 0)
	for _, metric_str := range metrics_str {
		fmt.Println(metric_str)
		var valuesItem ValuesStringItem
		split_metrics := strings.Split(metric_str, "' labels={")
		valuesItem.Metric = strings.ReplaceAll(split_metrics[0], " metric='", "")
		split_value := strings.Split(split_metrics[1], "} value=")
		valuesItem.Value = strings.ReplaceAll(split_value[1], "]", "")
		valuesItem.Labels = strings.Split(split_value[0], ",")

		resultsValues = append(resultsValues, valuesItem)
	}
	return resultsValues
}

type TemplateData struct {
	AlertName   string
	AlertURL    string
	AlertAtTime string
	AlertLabels string
	AlertValues string
	AlertKeepTime string
}

type ExtendedAlert struct {
	Status        string      `json:"status"`
	Labels        template.KV `json:"labels"`
	Annotations   template.KV `json:"annotations"`
	StartsAt      time.Time   `json:"startsAt"`
	EndsAt        time.Time   `json:"endsAt"`
	GeneratorURL  string      `json:"generatorURL"`
	Fingerprint   string      `json:"fingerprint"`
	SilenceURL    string      `json:"silenceURL"`
	DashboardURL  string      `json:"dashboardURL"`
	PanelURL      string      `json:"panelURL"`
	ValueString   string      `json:"valueString"`
	ImageURL      string      `json:"imageURL,omitempty"`
	EmbeddedImage string      `json:"embeddedImage,omitempty"`
}

type ExtendedAlerts []ExtendedAlert

type ExtendedData struct {
	Receiver string         `json:"receiver"`
	Status   string         `json:"status"`
	Alerts   ExtendedAlerts `json:"alerts"`

	GroupLabels       template.KV `json:"groupLabels"`
	CommonLabels      template.KV `json:"commonLabels"`
	CommonAnnotations template.KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}
