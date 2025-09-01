#!/bin/bash

touch .env
echo "OMDB_API_KEY=" >> .env
echo "SCRAPPER_PORT=8000" >> .env
echo "EXPOSED_PORT=8080" >> .env
