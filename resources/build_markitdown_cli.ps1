# Script to build markitdown-cli.exe using PyInstaller and a temporary virtual environment
# Requires Python to be installed and available on PATH.

$ResourcesDir = $PSScriptRoot
Push-Location $ResourcesDir

Write-Host "Creating Python virtual environment..." -ForegroundColor Cyan
if (Test-Path ".venv") {
    Remove-Item -Path ".venv" -Recurse -Force | Out-Null
}

python -m venv .venv
if (-not (Test-Path ".venv")) {
    Write-Error "Failed to create virtual environment."
    Pop-Location
    Exit 1
}

$VenvPython = Join-Path $ResourcesDir ".venv\Scripts\python.exe"
$VenvPip = Join-Path $ResourcesDir ".venv\Scripts\pip.exe"

Write-Host "Upgrading pip..." -ForegroundColor Cyan
Start-Process -FilePath $VenvPython -ArgumentList "-m pip install --upgrade pip" -Wait -NoNewWindow

Write-Host "Installing markitdown[all] and pyinstaller..." -ForegroundColor Cyan
Start-Process -FilePath $VenvPip -ArgumentList "install `"markitdown[all]`" pyinstaller" -Wait -NoNewWindow

Write-Host "Compiling markitdown_cli.py to executable..." -ForegroundColor Cyan
# Execute pyinstaller through venv's python module to ensure correct environment
$PyInstallerScript = Join-Path $ResourcesDir ".venv\Scripts\pyinstaller.exe"
$BuildArgs = "--onefile --console --name markitdown-cli --collect-data magika --collect-data markitdown `"$ResourcesDir\markitdown_cli.py`""
Start-Process -FilePath $PyInstallerScript -ArgumentList $BuildArgs -Wait -NoNewWindow

# Verify build output
$BuildOutput = Join-Path $ResourcesDir "dist\markitdown-cli.exe"
if (Test-Path $BuildOutput) {
    Write-Host "Compilation succeeded. Copying binary..." -ForegroundColor Green
    Copy-Item -Path $BuildOutput -Destination (Join-Path $ResourcesDir "markitdown-cli.exe") -Force
} else {
    Write-Error "Compilation failed. No executable was found in dist/."
    Pop-Location
    Exit 1
}

Write-Host "Cleaning up temporary build files..." -ForegroundColor Cyan
if (Test-Path "build") { Remove-Item -Path "build" -Recurse -Force | Out-Null }
if (Test-Path "dist") { Remove-Item -Path "dist" -Recurse -Force | Out-Null }
if (Test-Path "markitdown-cli.spec") { Remove-Item -Path "markitdown-cli.spec" -Force | Out-Null }

Write-Host "Build complete! markitdown-cli.exe is ready in $ResourcesDir" -ForegroundColor Green
Pop-Location
