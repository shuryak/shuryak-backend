package internal

import (
	"encoding/json"
	"log"
	"os"
)

var Configuration *ConfigType

type ConfigType struct {
	Debug   *ProfileType `json:"debug"`
	Release *ProfileType `json:"release"`
}

type ProfileType struct {
	ServerPort            *string `json:"server_port"`
	MongoConnectionString *string `json:"mongo_connection_string"`
}

func init() {
	file, err := os.Open("./configs/appsettings.json")
	if err != nil {
		log.Fatal("Invalid config path!\n\t>>> ", err)
	}
	decoder := json.NewDecoder(file)
	Configuration = new(ConfigType)
	err = decoder.Decode(&Configuration)
	if err != nil {
		log.Fatal("Bad config!")
	}
}
