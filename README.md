# Trailbrake API

## Description

Trailbrake API is an API for creating and retrieving ride data for the Trailbrake mobile application.

## Usage

See 'openapi.yaml' for the complete OpenAPI specifications.

## Setup

### Deploying to GKE

1. Create a Standard cluster: https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-zonal-cluster
   ```
   gcloud container clusters create {CLUSTER_NAME} --release-channel None --zone {COMPUTE_ZONE} --node-locations {COMPUTER_ZONE}
   ```
2. Build and push image to Google Artifact Registry: https://cloud.google.com/kubernetes-engine/docs/tutorials/hello-app; to ensure build compatibility, try building with Google Cloud Build: https://cloud.google.com/build/docs/building/build-containers
   ```
   gcloud builds submit --tag {LOCATION}-docker.pkg.dev/{PROJECT_ID}/{REPOSITORY}/{IMAGE_NAME}
   ```
3. Upload secrets to Google Secret Manager: https://cloud.google.com/secret-manager/docs/create-secret-quickstart
   ```
   echo -n "secret_data_content" | gcloud secrets create {SECRET_NAME} --replication-policy={automatic|user-managed} --data-file=-
   ```
4. Configure Workload Identity on GKE cluster to use IAM service accounts to access secrets by following instructions here: https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity

   ```
   gcloud container node-pools update {NODEPOOL_NAME} --cluster={CLUSTER_NAME} --region={COMPUTE_REGION} --workload-metadata=GKE_METADATA
   ```

   _This sets up `kubectl` command_

   ```
   gcloud container clusters get-credentials {CLUSTER_NAME} --region={COMPUTE_REGION}
   ```

   _This is the Google IAM Service Account (GSA) that will access secrets_

   ```
   gcloud iam service-accounts create {GSA_NAME} --project={GSA_PROJECT}
   ```

   _Grant secret accessor role to GSA_

   ```
   gcloud projects add-iam-policy-binding {PROJECT_ID} --member "serviceAccount:{GSA_NAME}@{GSA_PROJECT}.iam.gserviceaccount.com" --role "roles/secretmanager.secretAccessor"
   ```

   _Bind Kubernetes Service Account (KSA) to GSA to impersonate it_

   ```
   gcloud iam service-accounts add-iam-policy-binding {GSA_NAME}@{GSA_PROJECT}.iam.gserviceaccount.com --role roles/iam.workloadIdentityUser --member "serviceAccount:{PROJECT_ID}.svc.id.goog[{NAMESPACE}/{KSA_NAME}]"
   ```

   _Annotate KSA with GSA_

   ```
   kubectl annotate serviceaccount KSA_NAME --namespace {NAMESPACE} iam.gke.io/gcp-service-account={GSA_NAME}@{GSA_PROJECT}.iam.gserviceaccount.com
   ```

5. Install Secrets Store CSI Driver in cluster: https://secrets-store-csi-driver.sigs.k8s.io/getting-started/installation.html, or run `. k8s/install_secrets_csi_driver.sh`
6. Install GCP Secret Manager provider for Secret Store CSI Driver, and ensure IAM policy for service account is binded to secrets: https://github.com/GoogleCloudPlatform/secrets-store-csi-driver-provider-gcp, or run `. k8s/install_gcp_provider.sh`
7. Run `kubectl apply -f k8s/trailbrake-secrets.yaml`
8. Run `kubectl apply -f k8s/trailbrake-pod.yaml`
9. Run `kubectl apply -f k8s/trailbrake-services.yaml` to expose application as Services: https://cloud.google.com/kubernetes-engine/docs/how-to/exposing-apps

### Deploying to Cloud Run

1. `gcloud run services replace gcp/cloud-run.yaml`
