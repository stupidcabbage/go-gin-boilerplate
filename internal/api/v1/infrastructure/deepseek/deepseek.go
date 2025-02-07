package deepseek

import (
	"example.com/m/internal/config"
	"github.com/cohesion-org/deepseek-go"
)

var DeepseekClient *deepseek.Client

func InitClient() {
	DeepseekClient = deepseek.NewClient(config.Config.DeepseekApiKey)
}
