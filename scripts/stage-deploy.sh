#!/bin/bash
# LEP System - Stage Deploy (deploy no Cloud Run)
# Este script faz deploy da aplicação no Google Cloud Run para acesso online

set -e

echo "=== LEP System - Deploy STAGE no Google Cloud ==="
echo "• Deploy no Cloud Run"
echo "• Cloud SQL + GCS"
echo "• Acesso online via URL pública"
echo

# Verificar se gcloud está configurado
if ! command -v gcloud >/dev/null 2>&1; then
    echo "❌ Erro: gcloud CLI não encontrado. Por favor, instale o Google Cloud SDK."
    exit 1
fi

# Verificar autenticação
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" >/dev/null 2>&1; then
    echo "❌ Erro: Não autenticado no Google Cloud. Execute:"
    echo "   gcloud auth login"
    exit 1
fi

# Verificar projeto ativo
PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
if [ "$PROJECT_ID" != "leps-472702" ]; then
    echo "❌ Erro: Projeto Google Cloud incorreto. Execute:"
    echo "   gcloud config set project leps-472702"
    exit 1
fi

echo "✅ Google Cloud configurado: $PROJECT_ID"

# Verificar se terraform foi executado
if ! gcloud sql instances describe leps-postgres-stage --region=us-central1 >/dev/null 2>&1; then
    echo "❌ Erro: Infraestrutura não encontrada. Execute primeiro:"
    echo "   cd terraform && terraform init && terraform apply -var-file=../environments/gcp-stage.tfvars"
    exit 1
fi

echo "✅ Infraestrutura verificada"

# Configurações do deploy
REGION="us-central1"
SERVICE_NAME="lep-system"

echo "🚀 Iniciando deploy no Cloud Run..."
echo "   Região: $REGION"
echo "   Serviço: $SERVICE_NAME"
echo

# Fazer deploy usando gcloud build e deploy
ENVIRONMENT=stage gcloud run deploy $SERVICE_NAME \
    --source . \
    --region=$REGION \
    --platform=managed \
    --allow-unauthenticated \
    --memory=1Gi \
    --cpu=2 \
    --min-instances=1 \
    --max-instances=20 \
    --add-cloudsql-instances=leps-472702:us-central1:leps-postgres-stage \
    --service-account=lep-backend-sa@leps-472702.iam.gserviceaccount.com \
    --set-env-vars="ENVIRONMENT=stage,STORAGE_TYPE=gcs,BUCKET_NAME=leps-472702-lep-images-stage,BASE_URL=https://storage.googleapis.com/leps-472702-lep-images-stage,DB_USER=lep_user,DB_NAME=lep_database,INSTANCE_UNIX_SOCKET=/cloudsql/leps-472702:us-central1:leps-postgres-stage" \
    --set-secrets="DB_PASS=db-password-stage:latest,JWT_SECRET_PRIVATE_KEY=jwt-private-key-stage:latest,JWT_SECRET_PUBLIC_KEY=jwt-public-key-stage:latest"

# Obter URL do serviço
SERVICE_URL=$(gcloud run services describe $SERVICE_NAME --region=$REGION --format="value(status.url)")

echo
echo "=== Deploy STAGE Concluído com Sucesso! ==="
echo "🌐 URL do Serviço: $SERVICE_URL"
echo "🔍 Health Check: $SERVICE_URL/health"
echo "📊 Logs: gcloud logs read --limit=50 --format=json | jq -r '.textPayload'"
echo

# Testar conectividade
echo "🧪 Testando conectividade..."
if curl -f -s "$SERVICE_URL/ping" >/dev/null; then
    echo "✅ Serviço respondendo corretamente!"
else
    echo "⚠️  Serviço pode estar iniciando. Verifique os logs:"
    echo "   gcloud logs read --limit=10"
fi

echo
echo "Comandos úteis:"
echo "  gcloud run services logs read $SERVICE_NAME --region=$REGION  # Ver logs"
echo "  gcloud run services describe $SERVICE_NAME --region=$REGION   # Ver detalhes"
echo "  gcloud run services delete $SERVICE_NAME --region=$REGION     # Deletar serviço"
echo