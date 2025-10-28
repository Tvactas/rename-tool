# ===========================
# Build renamer.exe using Go
# ===========================
Write-Host "Building renamer.exe ..." -ForegroundColor Cyan

$build = Start-Process go -ArgumentList 'build', '-ldflags="-H windowsgui -s -w"', '-o', 'renamer.exe' -NoNewWindow -PassThru -Wait

if ($build.ExitCode -ne 0) {
    Write-Host "[ERROR] Go build failed!" -ForegroundColor Red
    exit $build.ExitCode
}

# ===========================
# Compress renamer.exe using UPX
# ===========================
Write-Host "Compressing with UPX ..." -ForegroundColor Cyan

$upxPath = "E:\upx-5.0.1-win64\upx-5.0.1-win64\upx.exe"
$upx = Start-Process $upxPath -ArgumentList '--best', '--lzma', 'renamer.exe' -NoNewWindow -PassThru -Wait

if ($upx.ExitCode -ne 0) {
    Write-Host "[ERROR] UPX compression failed!" -ForegroundColor Red
    exit $upx.ExitCode
}

Write-Host "Done!" -ForegroundColor Green
pause
