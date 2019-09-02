Param
(
    $OS=$null  
)

switch ($OS)
{
    "windows"
    {
        $env:GOOS = "windows"
    }
    "linux"
    {
        $env:GOOS = "linux"
    }
    "mac"
    {
        $env:GOOS = "darwin"
    }
}

$output = "ukfast"
if ($env:GOOS -eq "windows" -or ([string]::IsNullOrEmpty($env:GOOS) -and $Env:OS -eq "Windows_NT"))
{
    $output = $output+".exe"
}

$version = $(git describe --tags)
$builddate = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ss")
$env:GO111MODULE="on"

Write-Host -ForegroundColor Yellow -Object "Building $output with version $version and build date $builddate"
go build -o $output -ldflags "-s -X 'main.VERSION=$version' -X 'main.BUILDDATE=$builddate'"
$ec = $LASTEXITCODE

if (Test-Path -Path Env:\GOOS)
{
    Remove-Item -Path Env:\GOOS
}

exit $ec