#!/usr/bin/env bash
set -euo pipefail

PROJECT=dwk-gke-513717
CLUSTER=dwk-cluster
ZONE=europe-north1-b
REGION=europe-north1
REPOSITORY=my-repository
NAMESPACE=project
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT" --format='value(projectNumber)')
COMPUTE_SA="$PROJECT_NUMBER-compute@developer.gserviceaccount.com"
SA_NAME=storage-sa
SA_EMAIL="$SA_NAME@$PROJECT.iam.gserviceaccount.com"
BUCKET="gs://$PROJECT-backups"
BACKUP_KSA=project-backup-ksa

gcloud services enable \
	container.googleapis.com \
	artifactregistry.googleapis.com \
	storage.googleapis.com \
	--project "$PROJECT"

if ! gcloud container clusters describe "$CLUSTER" --zone "$ZONE" --project "$PROJECT" >/dev/null 2>&1; then
	gcloud container clusters create "$CLUSTER" \
		--project "$PROJECT" \
		--zone "$ZONE" \
		--cluster-version=1.36 \
		--disk-size=32 \
		--num-nodes=4 \
		--machine-type=e2-small \
		--gateway-api=standard \
		--workload-pool="$PROJECT".svc.id.goog
fi

gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE" --project "$PROJECT"

gcloud projects add-iam-policy-binding "$PROJECT" \
	--member="serviceAccount:$COMPUTE_SA" \
	--role=roles/artifactregistry.reader \
	--condition=None \
	--format=none

if ! gcloud artifacts repositories describe "$REPOSITORY" --location "$REGION" --project "$PROJECT" >/dev/null 2>&1; then
	gcloud artifacts repositories create "$REPOSITORY" \
		--location "$REGION" \
		--project "$PROJECT" \
		--repository-format=docker
fi

if ! gcloud storage buckets describe "$BUCKET" >/dev/null 2>&1; then
	gcloud storage buckets create "$BUCKET" \
		--project "$PROJECT" \
		--location "$REGION" \
		--uniform-bucket-level-access
fi

if ! gcloud iam service-accounts describe "$SA_EMAIL" --project "$PROJECT" >/dev/null 2>&1; then
	gcloud iam service-accounts create "$SA_NAME" \
		--project "$PROJECT" \
		--display-name="Pod Storage SA"
fi

storage_sa_roles=(
	"roles/storage.objectCreator"
	"roles/storage.objectViewer"
)
for role in "${storage_sa_roles[@]}"; do
	gcloud storage buckets add-iam-policy-binding "$BUCKET" \
		--member="serviceAccount:$SA_EMAIL" \
		--role="$role" \
		--format=none
done

gcloud iam service-accounts add-iam-policy-binding "$SA_EMAIL" \
	--project "$PROJECT" \
	--member="serviceAccount:$PROJECT.svc.id.goog[$NAMESPACE/$BACKUP_KSA]" \
	--role=roles/iam.workloadIdentityUser \
	--format=none
