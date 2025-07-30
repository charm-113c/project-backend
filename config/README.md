# The config package

This package/directory contains the needed utils to load the configurations. This package should contain two configuration files: config.yaml and server.env. These files ***must be*** .gitignored for security.
Alternatively, these files could be in some other directory, provided their correct path is inputted initially.

## Loading the configuration

The configurations are loaded using two imported modules: [cleanenv]("github.com/ilyakaznacheev/cleanenv") and [godotenv]("github.com/joho/godotenv"). godotenv is used in order to read .env files and so doing obtain the environmental variables, while cleanenv loads the environmental variables and the YAML config into the Golang `Config` struct.  
In theory, config files could be omitted entirely. In that case, cleanenv will rely on the default values provided in the `Config` struct, while the running mode of the server (development or production) can **and should** be set by running the appropriate `make` command from the Makefile.  
If config files are used, however, they must conform to the structure defined in the `Config` struct.

## Config.CfgPaths

The `Config` struct contains a field called `CfgPaths`. That field is not read by cleanenv or godotenv, but must instead be set manually by exporting the corresponding environmental variables. These variables must contain the path to the .yaml and .env configuration files, should they exist.
