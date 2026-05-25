
    

# Build [ Using Minikube ] 

* [ Tunneling the Port ] Keep It running in One port 
```sh
minikube tunnel
```

* In another terminal 
```sh
make k8-start
```

>> [Now We can see the build in http://localhost:9000/builds -> builds]

* To Compile the cli [ NOTE :: the binary is in ./bin/ folder ]
```sh
make build
```

# TO Build a Artifact 

How to run the code 
```sh
./bin/a2ctl <path/to/project>
```