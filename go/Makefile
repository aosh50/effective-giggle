APP_NAME=momenton
PORT=3333

build:
	docker build -t $(APP_NAME) .

run: 
	docker run --env-file=.env -p=$(PORT):$(PORT) --name="$(APP_NAME)" $(APP_NAME)