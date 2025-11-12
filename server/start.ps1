param([String]$testing)

$Env:DATABASE = "postgres"
$Env:DATABASE_IP = "localhost"
$Env:DATABASE_USER = "postgres"
$Env:DATABASE_PORT = 5432
$Env:DATABASE_PASSWORD = "postgres"
$Env:TLS_ENABLED = "true"
$Env:CERT_PATH = "../localhost.crt"
$Env:CERT_PRIVATE_KEY_PATH = "../localhost.key"
$Env:ALLOWED_ORIGIN = "*"

if($PSBoundParameters.ContainsKey('testing')){
    "[WARNING] Testing enabled, additional unsafe endpoints will be enabled"
    $Env:TESTING = "true"

} else{
    $Env:TESTING = "false"
}

go run .\main.go
