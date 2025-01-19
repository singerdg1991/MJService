package config

//type LoadTestCase struct {
//	Path           string
//	ConfigType     string
//	ConfigFileName string
//	Expected       error
//}

//var loadTestCases []LoadTestCase = []LoadTestCase{
//	{
//		Path:           "./configfolder",
//		ConfigType:     "xml",
//		ConfigFileName: "config.xml",
//		Expected:       errors.New("your environment variables info (path, configType, configFileName) is wrong"),
//	},
//	{
//		Path:           ".",
//		ConfigType:     "env",
//		ConfigFileName: ".config",
//		Expected:       errors.New("your .env config file not found"),
//	},
//	{
//		Path:           ".",
//		ConfigType:     "env",
//		ConfigFileName: ".env.sample",
//		Expected:       errors.New("your .env config file not found"),
//	},
//	{
//		Path:           ".",
//		ConfigType:     "env",
//		ConfigFileName: GetEnvPath(),
//		Expected:       nil,
//	},
//	{
//		Path:           "",
//		ConfigType:     "",
//		ConfigFileName: "",
//		Expected:       nil,
//	},
//}

//func TestLoad(t *testing.T) {
//	for _, loadTestCase := range loadTestCases {
//		err := Load(loadTestCase.Path, loadTestCase.ConfigType, loadTestCase.ConfigFileName)
//		if err != nil {
//			if loadTestCase.Expected == nil {
//				t.Error(err.Error())
//			}
//		} else {
//			if loadTestCase.Expected != nil {
//				t.Error(loadTestCase.Expected.Error())
//			}
//		}
//	}
//}
