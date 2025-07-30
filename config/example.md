# Configuration example

Note that .env files have priority over .yaml files,
meaning that if both are present, the .yaml config will be *overwritten* by the .env.

## Example .env file content

```
# Server variables
DEV_MODE=true
SRV_PORT=7777
SRV_HOST="0.0.0.0"
LOG_FILE="/var/log/junkyard.log"
CERT_FILE="/some/where/secure"
KEY_FILE="/some/where/secure"

# Database variables
DB_TYPE=postgres
DB_HOST=storage
DB_PORT=5432
DB_USER=postgres
DB_PWD=password
DB_NAME=database
DB_POOL_SIZE=20

# Router variables
READ_TIMEOUT=5s
WRITE_TIMEOUT=5s
# BEHIND_PROXY=false # A Fiber-specific setting

# Logto values
ENDPOINT="/logto/endpoint"
APP_ID=logtoAppID
APP_SECRET=logtoAppSecret
```

## Example .yaml file content

```
server:
  devMode: true
  port: 7777
  host: localhost
  logFile: "/var/log/junkyard.log"
  certFile: "/some/where/secure"
  keyFile: "/some/where/secure"

database:
  type: postgres
  host: storage
  port: 5432
  user: postgres
  password: password
  dbName: database
  poolSize: 20

router:
  readTimeout: 5s
  writeTimeout: 5s
  behindProxy: false

logto:
  endpoint: "/logto/endpoint"
  appID: logtoAppID
  appSecret: logtoAppSecret
```
