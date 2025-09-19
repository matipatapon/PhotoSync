param([String]$help, [String]$package="...", [String]$test=".")

if($PSBoundParameters.ContainsKey('help')){
    "test.ps1 | will run unit tests"
    "test.ps1 -package 'PACKAGE' | will run unit tests from given package"
    "test.ps1 -test 'regex' | will run unit tests matching given regex. Note it works best with --package"
} else{
    go test -v ./src/$package -run $test
    exit $LASTEXITCODE
}
