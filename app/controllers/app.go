package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/revel/revel"
	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type App struct {
	*revel.Controller
}

const (
	AnnotationKey = "config.fluentd.com/config"
)

func (c App) Validate() revel.Result {
	var request v1beta1.AdmissionReview
	var obj appsv1.Deployment
	c.Params.BindJSON(&request)

	rawObject := request.Request.Object.Raw

	err := json.Unmarshal(rawObject, &obj)
	if err != nil {
		msg := "Error occurred while deserializing request. Policy check will be skipped" + err.Error()
		fmt.Println(msg)
		response := v1beta1.AdmissionReview{
			Response: &v1beta1.AdmissionResponse{
				Allowed: false,
				Result:  &metav1.Status{Message: strings.TrimSpace(msg)},
			},
		}
		return c.RenderJSON(response)
	}

	config, ok := obj.Annotations[AnnotationKey]

	if ok {
		if err := validateFluentdConfig(config); err != nil {
			response := v1beta1.AdmissionReview{
				Response: &v1beta1.AdmissionResponse{
					Allowed: false,
					Result:  &metav1.Status{Message: err.Error()},
				},
			}
			return c.RenderJSON(response)
		}
	}

	response := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Allowed: true,
		},
	}
	return c.RenderJSON(response)

}

func validateFluentdConfig(config string) error {

	schema, err := jsonschema.Compile("schemas/template.json")
	if err != nil {
		return err
	}

	if err = schema.Validate(strings.NewReader(config)); err != nil {
		return err
	}

	return nil
}
