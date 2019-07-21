# Introduction

Scope of this software is just to play around with Kubernetes (Minikube). </br>
This software is exposing a set of rest-api to manage a collection of ToDo operations and is accessing mysql to store and read Todos. </br>
We will see how this app can be deployed locally, using just Docker containers to link one container to another or using kubernetes, minikube in this case. </br>
In the past I also used this app deployed on Pivotal CF.
https://github.com/DanielePalaia/cf-mysql-example </br>

## Datastore and rest api

The todos operations are saved in a mysql datastore defined in datastore.sql

```
CREATE TABLE ToDo (
	    ID int NOT NULL AUTO_INCREMENT,
	    Topic varchar(255),
	    Completed int,
	    Due varchar(255) DEFAULT '',
	    PRIMARY KEY (ID)
);
```

The software exposes these rest api which can be tested with curl

curl http://localhost:8080/todos
will get to the collection showing all the collection elements

this one will create a new element to the collection
curl -H "Content-Type: application/json" -d '{"Topic":"New TodoElem", "Completed":0}' -X POST http://localhost:8080/todos

this one will get an element:
curl http://localhost:8080/todos/1

this one will update an existing element of the collection
curl -H "Content-Type: application/json" -d '{"Id":0,"name":"New TodoElem Updated"}' -X PUT http://localhost:8080/todos

this one will delete a resource
curl -X DELETE http://localhost/todos/1

this one will delete all the collection
curl -X DELETE http://localhost/todos

## Testing the application locally:
Once built you can try the application locally: </br>

you need to create a mysql database as specified in datastore.sql file</br>

Then you can simply run the binary web-service-kubernetes

After it you can use curl to test the app </br>
You can test with curl the various rest api described before</br></br>

 
## Running the app on docker:
### Create a mysql docker like this: </br>
docker run -p 3306:3306 --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql/mysql-server:5.7 </br>
### Create database and new user
Enter inside the docker created and create the datastore as done locally </br>
docker exec -it some-mysql mysql -uroot -p</br>
</br>
After this create a new use and grant privileges on the database just created </br>
GRANT ALL PRIVILEGES ON *.* TO 'daniele'@'%' IDENTIFIED BY 'daniele' WITH GRANT OPTION; </br>
Do now a docker inspect some-mysql and get the ip of the docker image 
### Configure input properties 
Now collect all this info (username, password and ip and put it in the program configuration file ./conf
### Run the software in a docker container and link to mysql
A dockerfile is provided</br>
sudo  docker build -t web-service-kubernetes .</br>
docker run --publish 6060:8080 --name test --link some-mysql:mysql --rm web-service-kubernetes </br>
This will now listen on port 6060 use curl as done before to test it...</br>
### Test the rest api as before

 
## Running on kubernetes (minikube)

### Putting the docker image on dockerhub
I already created a dockerhub repository. In my case will be:</br>
https://cloud.docker.com/repository/registry-1.docker.io/danielepalaia/web-service-kubernetes
docker push danielepalaia/go-list:tagname</br>

### Install minikube
Minikube allows you to have and manage a local kubernetes cluster </br>
Follow this guide to install minikube on ubuntu </br>
https://linuxhint.com/install-minikube-ubuntu/</br>
</br>
Run minikube start and minikube dashboard to run the dashboard </br>
 ![Screenshot](./images/image1.png)

### Create a pod and a service for mysql
Follow this guide on how to create a mysql pod and service</br>
https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/

 ![Screenshot](./images/image2.png)

### Create a pod for this serviec web-service-kubernetes
I usually the minikube dashboard, you can go to new and specify as image danielepalaia/web-service-kubernetes

### Forward the port from pod locally
kubectl port-forward pod-name 8080:8080

### You can then use the rest api as before from your localhost

### Useful kubernetes command

1) Getting the shell from a container </br>

kubectl exec -it my-pod -- /bin/bash </br>

with docker instead: </br?

docker ps to take the id of the container

</br>
docker exec -it docker-id bash
</br></br?





 
 
