# LEP System - Interactive Multi-Environment Deployment Script (PowerShell)
# Supports: local-dev, gcp-dev, gcp-stage, gcp-prd

param(
    [string]$Environment = ""
)

# Configuration
$ProjectId = "leps-472702"
$ProjectName = "leps"
$Region = "us-central1"

# Global variables
$SelectedEnv = ""
$FailedCommands = @()
$RemainingCommands = @()
$CurrentStep = 0
$TotalSteps = 0

# Logging functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Write-Step {
    param([int]$Current, [int]$Total, [string]$Message)
    Write-Host "[STEP $Current/$Total] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Magenta
}

function Write-Command {
    param([string]$Command)
    Write-Host "[CMD] $Command" -ForegroundColor Cyan
}

# Show banner
function Show-Banner {
    Write-Host "=================================================================="  -ForegroundColor Magenta
    Write-Host "           LEP System - Interactive Deployment                    " -ForegroundColor Magenta
    Write-Host "=================================================================="  -ForegroundColor Magenta
    Write-Host ""
    Write-Host "Choose your deployment environment:"
    Write-Host ""
    Write-Host "1. 🏠 local-dev     - Local development with Docker Compose"
    Write-Host "2. ☁️  gcp-dev      - GCP minimal setup (testing only)"
    Write-Host "3. 🚀 gcp-stage     - GCP staging (production-like, no Twilio)"
    Write-Host "4. 🌟 gcp-prd       - GCP production (full features)"
    Write-Host ""
}

# Environment selection
function Select-Environment {
    if ($Environment -ne "") {
        switch ($Environment.ToLower()) {
            "local-dev" { $global:SelectedEnv = "local-dev"; return }
            "gcp-dev" { $global:SelectedEnv = "gcp-dev"; return }
            "gcp-stage" { $global:SelectedEnv = "gcp-stage"; return }
            "gcp-prd" { $global:SelectedEnv = "gcp-prd"; return }
            default {
                Write-Error "Invalid environment: $Environment"
                exit 1
            }
        }
    }

    while ($true) {
        $choice = Read-Host "Select environment (1-4)"
        switch ($choice) {
            "1" { $global:SelectedEnv = "local-dev"; break }
            "2" { $global:SelectedEnv = "gcp-dev"; break }
            "3" { $global:SelectedEnv = "gcp-stage"; break }
            "4" { $global:SelectedEnv = "gcp-prd"; break }
            default { Write-Error "Invalid choice. Please select 1-4."; continue }
        }
        break
    }

    Write-Info "Selected environment: $SelectedEnv"
}

# Check if command exists
function Test-Command {
    param([string]$Command)
    try {
        if (Get-Command $Command -ErrorAction SilentlyContinue) {
            return $true
        }
        return $false
    }
    catch {
        return $false
    }
}

# Dependency checks
function Test-Dependencies {
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Checking dependencies for $SelectedEnv"

    $missingDeps = @()

    switch -Wildcard ($SelectedEnv) {
        "local-dev" {
            if (-not (Test-Command "docker")) { $missingDeps += "docker" }
            if (-not (Test-Command "docker-compose")) { $missingDeps += "docker-compose" }
            if (-not (Test-Command "go")) { $missingDeps += "go" }
        }
        "gcp-*" {
            if (-not (Test-Command "gcloud")) { $missingDeps += "gcloud" }
            if (-not (Test-Command "terraform")) { $missingDeps += "terraform" }
            if (-not (Test-Command "docker")) { $missingDeps += "docker" }
            if (-not (Test-Command "go")) { $missingDeps += "go" }
        }
    }

    if ($missingDeps.Count -gt 0) {
        Write-Error "Missing required dependencies: $($missingDeps -join ', ')"
        Write-Host ""
        Write-Host "Installation instructions:"
        foreach ($dep in $missingDeps) {
            switch ($dep) {
                "docker" { Write-Host "  - Docker: https://docs.docker.com/desktop/windows/" }
                "docker-compose" { Write-Host "  - Docker Compose: Included with Docker Desktop" }
                "gcloud" { Write-Host "  - Google Cloud CLI: https://cloud.google.com/sdk/docs/install-windows" }
                "terraform" { Write-Host "  - Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli" }
                "go" { Write-Host "  - Go: https://golang.org/doc/install" }
            }
        }
        exit 1
    }

    Write-Success "All dependencies are installed"
}

# Environment variable validation
function Test-EnvironmentVars {
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Validating environment variables"

    $missingVars = @()
    $configFile = ""

    switch ($SelectedEnv) {
        "local-dev" {
            $configFile = "environments\local-dev.env"
            if (-not (Test-Path $configFile)) {
                $missingVars += "environments\local-dev.env file"
            }
        }
        "gcp-dev" {
            $configFile = "environments\gcp-dev.tfvars"
        }
        "gcp-stage" {
            $configFile = "environments\gcp-stage.tfvars"
            # Check for SMTP configuration
            if (Test-Path $configFile) {
                $content = Get-Content $configFile -Raw
                if ($content -notmatch 'smtp_username.*@') {
                    $missingVars += "SMTP_USERNAME in gcp-stage.tfvars"
                }
            }
        }
        "gcp-prd" {
            $configFile = "environments\gcp-prd.tfvars"
            # Check for Twilio configuration
            if (Test-Path $configFile) {
                $content = Get-Content $configFile -Raw
                if ($content -notmatch 'twilio_account_sid.*AC') {
                    $missingVars += "TWILIO_ACCOUNT_SID in gcp-prd.tfvars"
                }
                if ($content -notmatch 'smtp_username.*@') {
                    $missingVars += "SMTP_USERNAME in gcp-prd.tfvars"
                }
            }
        }
    }

    # Check config file exists
    if (-not (Test-Path $configFile)) {
        $missingVars += $configFile
    }

    # For GCP environments, check JWT keys
    if ($SelectedEnv -like "gcp-*") {
        if (-not (Test-Path "jwt_private_key.pem") -or -not (Test-Path "jwt_public_key.pem")) {
            $missingVars += "JWT key files (jwt_private_key.pem, jwt_public_key.pem)"
        }
    }

    if ($missingVars.Count -gt 0) {
        Write-Error "Missing required configuration: $($missingVars -join ', ')"
        Write-Host ""
        Write-Host "Required configurations:"
        foreach ($var in $missingVars) {
            Write-Host "  - $var"
        }
        exit 1
    }

    Write-Success "Environment validation passed"
}

# Execute command with error handling
function Invoke-DeployCommand {
    param([string]$Command, [string]$Description)

    Write-Command $Command

    try {
        if ($Command.StartsWith("docker-compose")) {
            $parts = $Command.Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
            $args = $parts[1..($parts.Length-1)]
            & docker-compose $args
        }
        elseif ($Command.StartsWith("gcloud")) {
            $parts = $Command.Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
            $args = $parts[1..($parts.Length-1)]
            & gcloud $args
        }
        elseif ($Command.StartsWith("terraform")) {
            $parts = $Command.Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
            $args = $parts[1..($parts.Length-1)]
            & terraform $args
        }
        elseif ($Command.StartsWith("docker")) {
            $parts = $Command.Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
            $args = $parts[1..($parts.Length-1)]
            & docker $args
        }
        else {
            Invoke-Expression $Command
        }

        if ($LASTEXITCODE -eq 0) {
            Write-Success "$Description completed"
            return $true
        }
        else {
            throw "Command failed with exit code $LASTEXITCODE"
        }
    }
    catch {
        Write-Error "$Description failed: $($_.Exception.Message)"
        $global:FailedCommands += $Command

        # Show remaining commands
        if ($global:RemainingCommands.Count -gt 0) {
            Write-Host ""
            Write-Warn "Remaining commands to complete deployment:"
            foreach ($remainingCmd in $global:RemainingCommands) {
                Write-Host "  - $remainingCmd"
            }
        }

        Write-Host ""
        Write-Error "Deployment failed at: $Description"
        Write-Error "Failed command: $Command"
        exit 1
    }
}

# Local development deployment
function Deploy-LocalDev {
    Write-Info "🏠 Deploying LOCAL DEVELOPMENT environment"

    $global:RemainingCommands = @(
        "docker-compose down",
        "docker-compose build",
        "docker-compose up -d",
        "docker-compose logs -f app"
    )

    $global:TotalSteps = 6

    Invoke-DeployCommand "docker-compose down -v" "Cleaning existing containers"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    Invoke-DeployCommand "Copy-Item environments\local-dev.env .env" "Setting up environment file"

    Invoke-DeployCommand "docker-compose build --no-cache app" "Building application container"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    Invoke-DeployCommand "docker-compose up -d postgres redis mailhog" "Starting infrastructure services"

    Write-Info "Waiting for database to be ready..."
    Start-Sleep 10

    Invoke-DeployCommand "docker-compose up -d app" "Starting application"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    Write-Success "🎉 Local development environment is ready!"
    Write-Host ""
    Write-Host "📍 Available services:"
    Write-Host "  - Application: http://localhost:8080"
    Write-Host "  - Health check: http://localhost:8080/health"
    Write-Host "  - Database (PgAdmin): http://localhost:5050 (admin@lep.local / admin123)"
    Write-Host "  - Email testing (MailHog): http://localhost:8025"
    Write-Host ""
    Write-Host "📋 Useful commands:"
    Write-Host "  - View logs: docker-compose logs -f app"
    Write-Host "  - Stop services: docker-compose down"
    Write-Host "  - Database shell: docker-compose exec postgres psql -U lep_user -d lep_database"
    Write-Host ""

    $showLogs = Read-Host "Show application logs? (y/n)"
    if ($showLogs -eq "y" -or $showLogs -eq "Y") {
        & docker-compose logs -f app
    }
}

function Deploy-GCP {
    $envSuffix = switch -Wildcard ($SelectedEnv) {
        "gcp-dev" { "dev" }
        "gcp-stage" { "staging" }
        "gcp-prd" { "prod" }
    }

    Write-Info "☁️ Deploying to GCP environment: $SelectedEnv"

    $global:RemainingCommands = @(
        "gcloud auth login",
        "gcloud config set project $ProjectId",
        "terraform init",
        "terraform plan -var-file=environments\$SelectedEnv.tfvars",
        "terraform apply -var-file=environments\$SelectedEnv.tfvars",
        "docker build",
        "docker push",
        "gcloud run deploy"
    )

    $global:TotalSteps = 10

    # Authenticate with GCP
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Authenticating with Google Cloud"
    try {
        $activeAccount = & gcloud auth list --filter=status:ACTIVE --format="value(account)" 2>$null
        if (-not $activeAccount) {
            Invoke-DeployCommand "gcloud auth login" "Google Cloud authentication"
            Invoke-DeployCommand "gcloud auth application-default login" "Application default credentials"
            Invoke-DeployCommand "gcloud auth configure-docker $Region-docker.pkg.dev" "Docker authentication"
        }
        else {
            Write-Success "Already authenticated with Google Cloud"
        }
    }
    catch {
        Invoke-DeployCommand "gcloud auth login" "Google Cloud authentication"
    }
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Set project
    Invoke-DeployCommand "gcloud config set project $ProjectId" "Setting GCP project"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Initialize Terraform
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Initializing Terraform"
    Invoke-DeployCommand "terraform init" "Terraform initialization"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Plan Terraform
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Planning infrastructure"
    Invoke-DeployCommand "terraform plan -var-file=environments\$SelectedEnv.tfvars -out=tfplan" "Terraform planning"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Confirm deployment
    Write-Host ""
    Write-Warn "Review the Terraform plan above."
    $confirm = Read-Host "Continue with deployment? (y/n)"
    if ($confirm -ne "y" -and $confirm -ne "Y") {
        Write-Warn "Deployment cancelled by user."
        exit 0
    }

    # Apply Terraform
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Creating infrastructure"
    Invoke-DeployCommand "terraform apply tfplan" "Infrastructure deployment"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Build and push Docker image
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Building Docker image"
    $imageUrl = "$Region-docker.pkg.dev/$ProjectId/lep-backend/lep-backend:latest"
    Invoke-DeployCommand "docker build -t $imageUrl ." "Docker image build"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    Write-Step (++$global:CurrentStep) $global:TotalSteps "Pushing Docker image"
    Invoke-DeployCommand "docker push $imageUrl" "Docker image push"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Deploy to Cloud Run
    Write-Step (++$global:CurrentStep) $global:TotalSteps "Deploying to Cloud Run"
    $serviceName = "$ProjectName-backend-$envSuffix"
    Invoke-DeployCommand "gcloud run deploy $serviceName --image=$imageUrl --region=$Region --platform=managed --allow-unauthenticated" "Cloud Run deployment"
    $global:RemainingCommands = $global:RemainingCommands[1..($global:RemainingCommands.Length-1)]

    # Get service URL
    try {
        $serviceUrl = & gcloud run services describe $serviceName --region=$Region --format="value(status.url)" 2>$null
    }
    catch {
        $serviceUrl = ""
    }

    Write-Success "🎉 GCP deployment completed!"
    Write-Host ""
    Write-Host "📍 Deployment information:"
    Write-Host "  - Environment: $SelectedEnv"
    Write-Host "  - Service: $serviceName"
    Write-Host "  - Region: $Region"
    if ($serviceUrl) {
        Write-Host "  - URL: $serviceUrl"
        Write-Host "  - Health check: $serviceUrl/health"
    }
    Write-Host ""

    # Health check
    if ($serviceUrl) {
        Write-Info "Performing health check..."
        Start-Sleep 10
        try {
            $response = Invoke-RestMethod -Uri "$serviceUrl/health" -Method Get -TimeoutSec 10
            Write-Success "Application is healthy!"
        }
        catch {
            Write-Warn "Health check failed. Service might still be starting up."
        }
    }
}

# Main execution
function Main {
    Show-Banner
    Select-Environment

    # Set total steps based on environment
    switch -Wildcard ($SelectedEnv) {
        "local-dev" { $global:TotalSteps = 6 }
        "gcp-*"     { $global:TotalSteps = 10 }
    }

    Test-Dependencies
    Test-EnvironmentVars

    Write-Host ""
    Write-Info "🚀 Starting deployment to $SelectedEnv"
    Write-Host ""

    switch -Wildcard ($SelectedEnv) {
        "local-dev" { Deploy-LocalDev }
        "gcp-*"     { Deploy-GCP }
        default     { Write-Error "Unsupported environment: $SelectedEnv"; exit 1 }
    }
}

# Execute main function
Main
