

start : 
	kubectl create -f namespace.yaml

	kubectl create -f shared-volumes-pod.yaml

	kubectl create -f builder-deployment.yaml
	kubectl create -f builder-service.yaml

	kubectl create -f consumer-deployment.yaml
	kubectl create -f consumer-service.yaml

stop :
	kubectl delete -f consumer-service.yaml
	kubectl delete -f builder-service.yaml


	kubectl delete -f consumer-deployment.yaml
	kubectl delete -f builder-deployment.yaml

	kubectl delete -f shared-volumes-pod.yaml
	kubectl delete -f namespace.yaml