# Weather-App

This repository hosts a simple web application for tracking weather using GO.


### Prerequisites

  Before running the application, ensure you have the following installed:

   - go version go1.18.1

### Installation

**1.** **Clone the repository and navigate to the project directory:**

```bash
git clone https://github.com/bassma-khaled14/Argocd-Demo-Project.git
cd Weather-App
```
**1.** **Install dependencies:**

```bash
go mid tidy
```
Set up Docker on your machine
```bash
sudo apt-get update
sudo apt install docker.io docker-compose docker-buildx
```
Test Docker 
```bash
sudo groupadd docker
docker run hello-world
```
Make sure that Docker service is enabled 
```bash
sudo systemctl status docker
```
If It's not enabled
```bash
sudo systemctl enable docker
```
### Usage
**1.** **Run the application:**
```bash
go run main
```
The application will be accessible at `http://localhost:8081`.
### Tests

 **2.** **Run test:**
````bash
go test
````
This ensures all functionalities are working correctly.

### Running the application in Docker
**1.** **Login With Docker**
````bash
docker login username
````
enter your username on dockerhub
   
**2.** **Pull Image From DockerHub**
    
````bash
docker pull bassma/weather-app:latest
````
**3.** **Run Docker Container**
````bash
docker run -p 5000:8080 --name appcontainer bassma/weather-app:latest
````
application will be accessible at http://127.0.0.1:5000

### Minikube Deployment

**1.** **install minikube** 
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
sudo apt install -y virtualbox virtualbox-ext-pack
sudo snap install kubectl --classic
kubectl config use-context minikube
```
**2.** **start minikube** 
```bash
minikube start
````
**3.** **Deploy**
```bash
# Deploy the application
kubectl apply -f deployment.yaml
# Create the service
kubectl apply -f service.yaml
````
**4.** **Accessing the Application**
```bash
minikube service <service> -n <namespace>
```
**5.** **To access a specific pod in a Kubernetes cluster**
```bash
kubectl get pod -n <namespace>
kubectl describe pod  -n <namespace>
```

### Install Argo CD on Kubernetes
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```
**1.** **Access the Argo CD UI**
Port-forward the Argo CD API server to your local machine:
```bash 
kubectl port-forward svc/argocd-server -n argocd 8080:80
```
This will forward Argo CD's API server to your local machine, and you can access the UI at http://localhost:8080.

**2.** **Log in to the Argo CD UI**
Default username: admin
Run this command to get the password:
```bash 
kubectl -n argocd get pods -l app.kubernetes.io/name=argocd-server -o name
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 --decode; echo
```
**3.** **Deploy Your App Using Argo CD**
``` bash
kubectl apply -f argo-app.yml
kubectl -n argocd argo sync my-app
```

![Argocd dashboard](assets/argocd-deployment.png)