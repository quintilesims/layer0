$ErrorActionPreference = 'Stop';


$packageName= 'layer0'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = 'https://github.com/quintilesims/layer0/releases/download/v0.10.7/Windows.zip'
$fileLocation = Join-Path $toolsDir 'l0.exe'
$unzipMethod = 'builtin'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  fileType      = 'EXE'
  url           = $url
  url64bit      = $url64
  file         = $fileLocation

  silentArgs    = "/qn /norestart /l*v `"$($env:TEMP)\$($packageName).$($env:chocolateyPackageVersion).MsiInstall.log`""
  validExitCodes= @(0, 3010, 1641)

  softwareName  = 'layer0*'
  checksum      = ''
  checksumType  = 'md5'
  checksum64    = ''
  checksumType64= 'md5'
}

Install-ChocolateyZipPackage @packageArgs
