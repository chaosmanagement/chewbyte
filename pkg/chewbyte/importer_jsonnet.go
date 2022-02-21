package chewbyte

import "github.com/google/go-jsonnet"

func importJsonnet(content, filenameHint string) (interface{}, error) {
	vm := jsonnet.MakeVM()

	jsonStr, err := vm.EvaluateAnonymousSnippet(filenameHint, content)
	if err != nil {
		return nil, err
	}

	return importJson(jsonStr)
}
