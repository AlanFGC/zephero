#!/bin/bash

remove_container() {
    local container_name="sqlite-container"
    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        printf "Removing existing container '%s'...\n" "$container_name"
        if ! docker rm -f "$container_name"; then
            printf "Failed to remove existing container '%s'.\n" "$container_name" >&2
            return 1
        fi
    fi
}

main() {
    printf "Creating world db...\n"
    if ! ./create_sample_db.sh world.db; then
        printf "Failed to create world db.\n" >&2
        return 1
    fi

    remove_container || return 1

    printf "Running SQLite container...\n"
    if ! docker run -it --name sqlite-container -v "$(pwd)/sample.db:/sample.db" alpine-sqlite; then
        printf "Failed to run SQLite container.\n" >&2
        return 1
    fi
}

main "$@"
