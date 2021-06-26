ifeq ($(env), prod)
	url = http://35.202.163.76:1323
else
	url = http://localhost:1323
endif

up:
	docker-compose up -d

down:
	docker-compose down

exec-db:
	docker exec -it app_db_1 bash

ps:
	docker ps

logs:
	docker logs free_go_1 -t

test-post-failed:
	curl $(url)/recipes \
	-H 'Content-Type:application/json' \
	-d "{"cost": 1000}"

test-post-success:
	curl $(url)/recipes \
	-X POST \
	-H 'Content-Type:application/json' \
	-d '{"cost":1000, "making_time":"20", "serves":"3", "ingredients":"eggs", "title":"tomato"}'

test-patch-all:
	curl $(url)/recipes/${id} \
	-X PATCH \
	-H 'Content-Type:application/json' \
	-d '{"cost":5000, "making_time":"30", "serves":"3", "ingredients":"eggs", "title":"tomato"}'	

test-patch:
	curl $(url)/recipes/${id} \
	-X PATCH \
	-H 'Content-Type:application/json' \
	-d '{"ingredients":"eggs, mayo, water", "title":"mashroom", "making_time":"100"}'

test-delete:
	curl -X DELETE $(url)/recipes/${id}

restart:
	docker-compose restart
