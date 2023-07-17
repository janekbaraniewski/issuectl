package issuectl

func ListRepositories() error {
	config := LoadConfig()

	Log.Infof("%v", config.Repositories)

	return nil
}

func GetRepository(name RepoConfigName) *RepoConfig {
	config := LoadConfig()

	for _, rc := range config.Repositories {
		if rc.Name == name {
			return &rc
		}
	}

	return nil
}

func AddRepository(repoConfig *RepoConfig) error {
	config := LoadConfig()

	config.Repositories = append(config.Repositories, *repoConfig)

	if err := config.Save(); err != nil {
		Log.Infof("ERROR!!!! %v", err)
		return err
	}

	return nil
}
