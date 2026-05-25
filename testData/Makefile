
help : 
	@echo "Run : Start the docker containers"
	@echo "Stop : Stop the docker containers"
	@echo "clean : Remove all the generated files and folders"

Run : 
	mkdir -p temps/storage
	mkdir -p temps/artifacts
	docker compose -f a2fManager.yml up # -d 

Stop : 
	docker compose -f a2fManager.yml down

clean : 
	rm -r temps/
	
.DEFAULT : help