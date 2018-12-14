echo off
set name=%1
docker build -t %1 -f DockerFile .