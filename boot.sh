#!/usr/bin/env bash
set -euo pipefail

PROJECT=dwk-gke-513717
CLUSTER=dwk-cluster
ZONE=europe-north1-b
REGION=europe-north1
REPOSITORY=my-repository

if ! gcloud container clusters describe "$CLUSTER" --zone "$ZONE" --project "$PROJECT" >/dev/null 2>&1; then
  gcloud container clusters create "$CLUSTER" \
    --project "$PROJECT" \
    --zone "$ZONE" \
    --cluster-version=1.36 \
    --disk-size=32 \
    --num-nodes=4 \
    --machine-type=e2-small \
    --gateway-api=standard
fi

gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE" --project "$PROJECT"

PROJECT_NUMBER=$(gcloud projects describe "$PROJECT" --format='value(projectNumber)')
gcloud projects add-iam-policy-binding "$PROJECT" \
  --member="serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role=roles/artifactregistry.reader \
  --condition=None \
  --format=none

if ! gcloud artifacts repositories describe "$REPOSITORY" --location "$REGION" --project "$PROJECT" >/dev/null 2>&1; then
  gcloud artifacts repositories create "$REPOSITORY" --location "$REGION" --project "$PROJECT" --repository-format=docker
fi

kubectl get nodes
