# LEP Backend Deployment Script for Windows PowerShell
# This script automates the deployment process to Google Cloud Platform

param(
    [string]$Environment = "dev",
    [switch]$SkipBuild,
    [switch]$SkipTerraform,
    [switch]$DryRun,
    [switch]$Help
)

# Colors for output
$Colors = @{
    Red = "Red"
    Green = "Green"
    Yellow = "Yellow"
    Blue = "Blue"
    White = "White"
}

function Write-LogInfo {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Blue
}

function Write-LogSuccess {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Green
}

function Write-LogWarning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Yellow
}

function Write-LogError {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Red
}

function Show-Usage {
    @"
Usage: .\deploy.ps1 [OPTIONS]

Deploy LEP Backend to Google Cloud Platform

OPTIONS:
    -Environment ENV     Environment to deploy (dev, staging, prod) [default: dev]
    -SkipBuild          Skip Docker image build
    -SkipTerraform      Skip Terraform infrastructure deployment
    -DryRun             Show what would be deployed without executing
    -Help               Show this help message

EXAMPLES:
    .\deploy.ps1 -Environment prod
    .\deploy.ps1 -SkipBuild -Environment staging
    .\deploy.ps1 -DryRun -Environment prod

PREREQUISITES:
    - gcloud CLI installed and authenticated
    - Docker Desktop installed and running
    - Terraform installed
    - terraform.tfvars file configured
    - JWT keys generated and configured

"@
}

function Test-Prerequisites {
    Write-LogInfo "Checking prerequisites..."

    # Check gcloud
    if (-not (Get-Command gcloud -ErrorAction SilentlyContinue)) {
        Write-LogError "gcloud CLI is not installed"
        exit 1
    }

    # Check Docker
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-LogError "Docker is not installed"
        exit 1
    }

    # Check Terraform
    if (-not (Get-Command terraform -ErrorAction SilentlyContinue)) {
        Write-LogError "Terraform is not installed"
        exit 1
    }

    # Check terraform.tfvars
    if (-not (Test-Path "terraform.tfvars")) {
        Write-LogError "terraform.tfvars file not found"
        Write-LogInfo "Copy terraform.tfvars.example to terraform.tfvars and configure it"
        exit 1
    }

    # Check if user is authenticated with gcloud
    $authCheck = gcloud auth list --filter=status:ACTIVE --format="value(account)" 2>$null
    if (-not $authCheck) {
        Write-LogError "Not authenticated with gcloud"
        Write-LogInfo "Run: gcloud auth login"
        exit 1
    }

    Write-LogSuccess "All prerequisites checked"
}

function Get-TerraformOutput {
    param([string]$OutputName)

    try {
        $output = terraform output -raw $OutputName 2>$null
        return $output
    }
    catch {
        return ""
    }
}

function Deploy-Infrastructure {
    if ($SkipTerraform) {
        Write-LogInfo "Skipping Terraform infrastructure deployment"
        return
    }

    Write-LogInfo "Deploying infrastructure with Terraform..."

    if ($DryRun) {
        Write-LogInfo "DRY RUN: Would run terraform plan and apply"
        return
    }

    # Initialize Terraform
    terraform init

    # Plan deployment
    terraform plan -var="environment=$Environment" -out=tfplan

    # Apply deployment
    Write-LogInfo "Applying Terraform configuration..."
    terraform apply tfplan

    # Clean up plan file
    Remove-Item tfplan -ErrorAction SilentlyContinue

    Write-LogSuccess "Infrastructure deployed successfully"
}

function Build-AndPushImage {
    if ($SkipBuild) {
        Write-LogInfo "Skipping Docker image build"
        return
    }

    Write-LogInfo "Building and pushing Docker image..."

    # Get repository URL from Terraform output
    $dockerRepoUrl = Get-TerraformOutput "docker_repository_url"

    if (-not $dockerRepoUrl) {
        Write-LogError "Could not get Docker repository URL from Terraform output"
        Write-LogInfo "Make sure infrastructure is deployed first"
        exit 1
    }

    $imageTag = "$dockerRepoUrl/lep-backend:latest"
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $imageTagEnv = "$dockerRepoUrl/lep-backend:$Environment-$timestamp"

    if ($DryRun) {
        Write-LogInfo "DRY RUN: Would build and push:"
        Write-LogInfo "  - $imageTag"
        Write-LogInfo "  - $imageTagEnv"
        return
    }

    # Configure Docker to authenticate with GCP
    Write-LogInfo "Configuring Docker authentication..."
    $registryHost = ($dockerRepoUrl -split '/')[0]
    gcloud auth configure-docker $registryHost --quiet

    # Build image
    Write-LogInfo "Building Docker image..."
    docker build -f Dockerfile.prod -t $imageTag -t $imageTagEnv .

    # Push images
    Write-LogInfo "Pushing Docker images..."
    docker push $imageTag
    docker push $imageTagEnv

    Write-LogSuccess "Docker images built and pushed successfully"
    Write-LogInfo "Latest image: $imageTag"
    Write-LogInfo "Tagged image: $imageTagEnv"
}

function Deploy-Service {
    Write-LogInfo "Deploying Cloud Run service..."

    $serviceName = Get-TerraformOutput "service_name"
    $region = Get-TerraformOutput "region"
    $imageUrl = Get-TerraformOutput "docker_image_url"

    if (-not $serviceName -or -not $region -or -not $imageUrl) {
        Write-LogError "Could not get service information from Terraform output"
        exit 1
    }

    if ($DryRun) {
        Write-LogInfo "DRY RUN: Would deploy Cloud Run service:"
        Write-LogInfo "  Service: $serviceName"
        Write-LogInfo "  Region: $region"
        Write-LogInfo "  Image: $imageUrl"
        return
    }

    # Deploy to Cloud Run
    gcloud run deploy $serviceName `
        --image=$imageUrl `
        --region=$region `
        --platform=managed `
        --quiet

    Write-LogSuccess "Cloud Run service deployed successfully"

    # Get service URL
    $serviceUrl = Get-TerraformOutput "service_url"
    Write-LogSuccess "Service available at: $serviceUrl"
}

function Test-Deployment {
    Write-LogInfo "Testing deployment..."

    $serviceUrl = Get-TerraformOutput "service_url"

    if (-not $serviceUrl) {
        Write-LogWarning "Could not get service URL for testing"
        return
    }

    if ($DryRun) {
        Write-LogInfo "DRY RUN: Would test endpoints:"
        Write-LogInfo "  - $serviceUrl/health"
        Write-LogInfo "  - $serviceUrl/ping"
        return
    }

    # Test health endpoint
    Write-LogInfo "Testing health endpoint..."
    try {
        $response = Invoke-WebRequest -Uri "$serviceUrl/health" -Method Get -TimeoutSec 10
        if ($response.StatusCode -eq 200) {
            Write-LogSuccess "Health check passed"
        }
        else {
            Write-LogWarning "Health check failed with status: $($response.StatusCode)"
        }
    }
    catch {
        Write-LogWarning "Health check failed: $($_.Exception.Message)"
    }

    # Test ping endpoint
    Write-LogInfo "Testing ping endpoint..."
    try {
        $response = Invoke-WebRequest -Uri "$serviceUrl/ping" -Method Get -TimeoutSec 10
        if ($response.StatusCode -eq 200) {
            Write-LogSuccess "Ping check passed"
        }
        else {
            Write-LogWarning "Ping check failed with status: $($response.StatusCode)"
        }
    }
    catch {
        Write-LogWarning "Ping check failed: $($_.Exception.Message)"
    }
}

# Main execution
if ($Help) {
    Show-Usage
    exit 0
}

# Validate environment
if ($Environment -notin @("dev", "staging", "prod")) {
    Write-LogError "Invalid environment: $Environment"
    Write-LogInfo "Valid environments: dev, staging, prod"
    exit 1
}

Write-LogInfo "Starting deployment for environment: $Environment"

if ($DryRun) {
    Write-LogWarning "DRY RUN MODE - No actual changes will be made"
}

Test-Prerequisites
Deploy-Infrastructure
Build-AndPushImage
Deploy-Service
Test-Deployment

Write-LogSuccess "Deployment completed successfully!"

if (-not $DryRun) {
    $serviceUrl = Get-TerraformOutput "service_url"
    Write-LogInfo "Service URL: $serviceUrl"
}