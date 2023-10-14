
# This will run both the terraform init and terraform validate commands. The init
# command will download the libraries and executables needed by the scripts that have been written to
# build our environment. The validate command will also validate our scripts are correct.
ready: init validate

# Use this command to easily run all of the infrastructure deployment steps
# There are costs associated with this demo; Remember to tear everything down!
deploy: plan apply

init:
	terraform init -upgrade

validate:
	terraform validate

plan:
	terraform plan -out app.tfplan

apply:
	terraform apply app.tfplan
	@echo !!!
	@echo !!!
	@echo !!! There are cost associated with running this demo! Remember to tear everything down!
	@echo !!!
	@echo !!! Use 'make destroy' OR use 'terraform destroy' to tear everything down.
	@echo !!!
	@echo !!!

destroy:
	terraform destroy -auto-approve
