#! /usr/bin/env bash

# for offset in 0 100 200 300 400 500 600 700 800 900 1000 1100 1200 1300 1400 1500; do
#   echo gorun main.go list --p "__Liked Songs 1536" --offset $offset >>likes.txt
# done

for i in {0..15}; do
  godotenv -f .env go run main.go list --p "__Liked Songs 1536" --offset $((i * 100)) >>likes.txt
done
