#!/bin/bash

init() {
	echo "Builiding the containers...."
    docker-compose up --build
}

case "$1" in
    init)
		init
		;;
    up)
		docker-compose up
		;;
    db)
        docker-compose exec user_service php setup.php
        docker-compose exec subscription_service php setup.php
		;;
    clean)
		docker-compose down -v
		;;
    test)
		echo "Testing {$2}....."
        docker-compose exec $2 php setup.php --mode=testing
		docker-compose exec $2 go clean -testcache
		docker-compose exec $2 go test -v ./tests/
		;;
	*)
		echo "Invalid option\n"
		;;
esac

