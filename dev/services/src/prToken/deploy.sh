microk8s.kubectl create -f prtoken-deployment.yaml
microk8s.kubectl create -f prtoken-service.yaml
microk8s.kubectl create -f roach-ui-claim0-persistentvolumeclaim.yaml 
microk8s.kubectl create -f roach-ui-deployment.yaml 
microk8s.kubectl create -f roach-ui-service.yaml
