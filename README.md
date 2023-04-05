# Trailbrake API

## Description

Trailbrake API is an API for creating and retrieving ride data for the Trailbrake mobile application.

## Usage

See 'openapi.yaml' for the complete OpenAPI specifications.

## Setup

Deploying to GKE:

1. Create a Standard cluster: https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-zonal-cluster
2. Build and push image to Google Artifact Registry: https://cloud.google.com/kubernetes-engine/docs/tutorials/hello-app
3. Upload secrets to Google Secret Manager: https://cloud.google.com/secret-manager/docs/create-secret-quickstart
4. Configure Workload Identity on GKE cluster to use IAM service accounts to access secrets byfollowing instructions here: https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity
5. Install Secrets Store CSI Driver in cluster: https://secrets-store-csi-driver.sigs.k8s.io/getting-started/installation.html, or run `. k8s/install_secrets_csi_driver.sh`
6. Install GCP Secret Manager provider for Secret Store CSI Driver, and ensure IAM policy for service account is binded to secrets: https://github.com/GoogleCloudPlatform/secrets-store-csi-driver-provider-gcp, or run `. k8s/install_gcp_provider.sh`
7. Run `kubectl apply -f k8s/trailbrake-secrets.yaml`
8. Run `kubectl apply -f k8s/trailbrake-pod.yaml`
9. Run `kubectl apply -f k8s/trailbrake-services.yaml` to expose application as Services: https://cloud.google.com/kubernetes-engine/docs/how-to/exposing-apps

Deploying to Cloud Run:

1. 
