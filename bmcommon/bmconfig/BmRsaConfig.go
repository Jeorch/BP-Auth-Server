package bmconfig

import "github.com/alfredyang1986/blackmirror/bmconfighandle"

type BmRsaConfig struct {
	Company string
	Date    string
}

func (br *BmRsaConfig) GenerateConfig() {
	configPath := "resource/rsaconfig.json"
	profileItems := bmconfig.BMGetConfigMap(configPath)

	br.Company = profileItems["Company"].(string)
	br.Date = profileItems["Date"].(string)
}
