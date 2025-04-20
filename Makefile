migrate:
	@docker run --rm -v ${PWD}/deployment/migration:/migrations \
	 migrate/migrate \
	 -path='/migrations/' \
	 -database=${DATABASE_URL} up