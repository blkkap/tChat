package config

import(
	"os"
	"path/filepath"
	"log"
)

func ConfigDir() string{
	dir,err := os.UserConfigDir()
	if err != nil{
		log.Fatal(err)
	}
	path := filepath.Join(dir, "tchat")
	
	err = os.MkdirAll(path, 0755)
	if err != nil{
		log.Fatal(err)
	}
	return path
}
//edits should go here if user naming is allowed
func ServerConfigPath() string {
	return filepath.Join(ConfigDir(),"server", "serverConfig.json")
}
func ClientConfigPath() string {
	return filepath.Join(ConfigDir(), "client", "clientConfig.json")
}
