$ErrorActionPreference = 'Stop';


$packageName= 'layer0'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = 'https://github.com/quintilesims/layer0/releases/download/v0.10.7/Windows.zip'
$unzipMethod = 'builtin'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  fileType      = 'EXE'
  url           = $url

  softwareName  = 'layer0*'
  checksum      = '2B86012BE3B71189E012DDCA32D39BE3'
  checksumType  = 'md5'
}

Install-ChocolateyZipPackage @packageArgs