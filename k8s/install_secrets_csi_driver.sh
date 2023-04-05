SECRETS_CSI_DRIVER_REPO=https://raw.githubusercontent.com/kubernetes-sigs/secrets-store-csi-driver/main
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/rbac-secretproviderclass.yaml
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/csidriver.yaml
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/secrets-store.csi.x-k8s.io_secretproviderclasses.yaml
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/secrets-store.csi.x-k8s.io_secretproviderclasspodstatuses.yaml
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/secrets-store-csi-driver.yaml

# If using the driver to sync secrets-store content as Kubernetes Secrets, deploy the additional RBAC permissions
# required to enable this feature
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/rbac-secretprovidersyncing.yaml

# If using the secret rotation feature, deploy the additional RBAC permissions
# required to enable this feature
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/rbac-secretproviderrotation.yaml

# If using the CSI Driver token requests feature (https://kubernetes-csi.github.io/docs/token-requests.html) to use
# pod/workload identity to request a token and use with providers
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/rbac-secretprovidertokenrequest.yaml

# [OPTIONAL] To deploy driver on windows nodes
kubectl apply -f $SECRETS_CSI_DRIVER_REPO/deploy/secrets-store-csi-driver-windows.yaml