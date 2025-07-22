package main

import (
	cfg "github.com/conductorone/baton-bamboohr/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("bamboohr", cfg.Config)
}
