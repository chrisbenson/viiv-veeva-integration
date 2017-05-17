package veeva

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/chrisbenson/easyaws/pkg/easyaws"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

var settings Settings

var conn string

type Settings struct {
	hour		string
	minute		string
	Luckie struct {
		DBUser string
		DBPassword string
		DBServer string
		DBName string
		DBDriver string
	}
	GSK struct {
		Username string
		Password string
		IdsURL string
		IdsNonce string
		CdiURL string
		CdiAuthToken string
		CdiNonce string
	}
}

func LoadConfig(awsProfile string) error {
	settings = Settings{}
	var awsSession *session.Session
	if awsProfile == "" {
		awsSession = easyaws.SessionFromEnvVars()
	} else {
		awsSession = easyaws.SessionFromProfile(awsProfile)
	}
	if awsSession == nil {
		return errors.New("Did not secure a valid AWS Session to get the config file from S3.")
	}
	key := "veeva.toml"
	keys := []string{key}
	tomlByteMap, err := easyaws.BytesFromS3("luckie-veeva", keys, awsSession)
	if err != nil {
		return errors.Wrap(err, "Failed to download TOML configuration file from S3 bucket.")
	}
	tomlBytes := tomlByteMap[key]
	if tomlBytes == nil {
		return errors.New("No bytes in tomlBytes for the key: " + key)
	}
	config, err := toml.Load(string(tomlBytes))
	if err != nil {
		return errors.Wrap(err, "Failed to load TOML from configuration file bytes received from S3.")
	} else {
		settings.hour			= config.Get("hour").(string)
		settings.minute			= config.Get("minute").(string)
		settings.Luckie.DBUser 		= config.Get("luckie.username").(string)
		settings.Luckie.DBPassword 	= config.Get("luckie.password").(string)
		settings.Luckie.DBServer 	= config.Get("luckie.server").(string)
		settings.Luckie.DBName 		= config.Get("luckie.name").(string)
		settings.Luckie.DBDriver 	= config.Get("luckie.driver").(string)
		settings.GSK.Username 		= config.Get("gsk.username").(string)
		settings.GSK.Password 		= config.Get("gsk.password").(string)
		settings.GSK.IdsURL 		= config.Get("gsk.ids_url").(string)
		settings.GSK.IdsNonce 		= config.Get("gsk.ids_nonce").(string)
		settings.GSK.CdiURL 		= config.Get("gsk.cdi_url").(string)
		settings.GSK.CdiAuthToken 	= config.Get("gsk.cdi_authtoken").(string)
		settings.GSK.CdiNonce 		= config.Get("gsk.cdi_nonce").(string)
		conn = "server=" + settings.Luckie.DBServer + ";user id=" + settings.Luckie.DBUser + ";password=" + settings.Luckie.DBPassword + ";database=" + settings.Luckie.DBName
	}
	return nil
}