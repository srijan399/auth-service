docker-up:
	docker compose up -d

docker-down:
	docker compose down

# delete volume
docker-no-volume:
	docker compose down -v

# docker check volume 
docker-volume-check:
	docker volume ls

# connect to postgres
postgres-connect:
	docker exec -it postgres-db psql -U srijan -d goauth