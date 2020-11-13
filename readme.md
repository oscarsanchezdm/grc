
# GRC

**G**o-Lang **R**est API server with a **C**RUD database

A project that uses Kubernetes for deploying a REST API server which communicates with a Mongo DB.

### About the project

The goal of the project is learning how to implement a REST-API server with Go, linked with a Mongo database, and using Kubernetes for the deployment. In order to simulate the Kubernetes' pods, minikube was used.

This project was developed in three days, starting with no kind of knowledge of Kubernetes, Go servers nor MongoDB-Go integration.

### Database
The database stores book entries. Every entry consists of an identifier (the ISBN), a title, an author and a boolean about the availability. Please note that this is a learning project, thus it is not intended to simulate or recreate real situations. The identifier might not be unique, as a random function is used to generate an integer.

### Preparing the environment
Before cloning this repo, you will need to [install Docker in your machine](https://docs.docker.com/desktop/). Intel Virtualization Technology is recommended to be enabled in your BIOS. Also, a Docker Hub account had to be created to push the containers images into its cloud, so make sure Docker is logged in in your account.

As the server is implemented with Go, you will need to have it installed. [Install instructions](https://golang.org/doc/install).

[kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) allows the user to communicate with the Kubernetes cluster, in this case, [minikube](https://kubernetes.io/es/docs/tasks/tools/install-minikube/).

Once all the tools have been installed, start Docker and the Kubernetes cluster

```
minikube start
```
In macOS, you will need to use the HyperKit driver for enabling the ingress, so instead of using the previous command, use
```
minikube start --driver=hyperkit
```
Once it has started, enable the ingress addon by
```
minikube addons enable ingress
```
We will generate a namespace in kubectl called grc to make things easier
```
kubectl create ns grc
```
### Deploying the server
The first step will be cloning the repository and making the shell point to the project's directory
```
git clone https://github.com/oscarsanchezdm/grc
cd grc
```
Before deploying the server, Go will need to install the required packages. To do so, simply run
```
go get
```
which will simply read go.mod file and download the packages automatically.

Once you have the required packages, it's time to compile using make
```
make
```
Once it has been compiled, docker will need to build the image before pushing it to Docker Hub. Note that these files are thought to be pushed in my Docker account, so for pushing it in another one you will need to change the project name in the following instructions
```
docker build -t oscarsanchezdm/grc
docker push oscarsanchezdm/grc
```

As we have the project image uploaded in Docker Hub server, we will make minikube to use it for the deployment. Minikube will deploy the API server and a Mongo server. For doing that, type
```
kubectl apply -n grc -f kubernetes/mongodb.yml
kubectl apply -n grc -f kubernetes/grc.yml
kubectl apply -n grc -f kubernetes/ingress.yml
```
[Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) will expose the REST API server to the outside. We will get the external IP address of the Ingress point using the following command
```
kubectl ingress -n grc
```
The output address will allow us to access to the API endpoints.

### Server endpoints
The server has several endpoints, which are defined below.  The curl examples use an IP which might not be the same as yours. Other methods to test the API such as PostMan can be used.
* **[POST/PUT] /grc** Creates a new entry in the database. A value will be sent in the request (see the examples section) containing the title and the author of the book. The ISBN will be randomly generated.
`curl -X POST -d "title=Un home cau&author=Jordi Basté" 192.168.64.4/grc/`
* **[GET] /grc** This endpoints returns a JSON containing all the book entries stored in the database.
`curl -X GET 192.168.64.4/grc/`
* **[GET] /grc/SearchByISBN/{isbn}** This endpoints returns a JSON containing all the book entries that have the specified ISBN.
`curl -X GET 192.168.64.4/grc/SearchByISBN/5765553123`
* **[GET] /grc/SearchByTitle/{title}** This endpoints returns a JSON containing all the book entries that have the specified book title.
`curl -X GET 192.168.64.4/grc/SarchByTitle/"Un home cau"`
* **[PATCH] /grc/{isbn}** This endpoints will change the availability boolean from false to true or from true to false.
 `curl -X PATCH 192.168.64.4/grc/5765553123`
* **[DELETE] /grc/{isbn}** This endpoints will allow to delete the book entries that have the specified ISBN.
`curl -X DELETE 192.168.64.4/grc/9272555810`

<<<<<<< HEAD
### Recompiling the API server while it is running
If the server needs to be recompiled, we will need to stop all the replicas and start them again, when the Docker Hub image has been updated.
=======
### Re-compiling the API server while it is running
If the server needs to be re-compiled, we will need to stop all the replicas and start them again, when the Docker Hub image has been updated.
>>>>>>> 3dc6ba26f77823921dc2fdfe3e63cf75dc8fd9f5
```
make
docker build -t oscarsanchezdm/grc
docker push oscarsanchezdm/grc
kubectl scale -n grc --replicas=0 delopyment grc
kubectl scale -n grc --replicas=1 delopyment grc
```
