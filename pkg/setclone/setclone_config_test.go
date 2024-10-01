package setclone

import (
	"gopkg.in/yaml.v2"
	"testing"
)

//func (receiver ) name()  {
//
//}

func TestReadConfig(t *testing.T) {
	// use setclone.yaml
	// call setclond_config.ReadConfig

	config := &setcloneConfig{}
	//setcloneConfig.ReadConfig("../../../test")
	//config.ReadConfig("../../../test")
	config.ReadConfig("./")
	t.Log(yaml.Marshal(config))
	// print unmarshal yaml format
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(yamlData))
}
