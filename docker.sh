docker build -t auth-api .
docker run -d --name auth-api \
  --add-host=host.docker.internal:host-gateway \
  --env-file docker.env \
  -p 8080:8080 \
  auth-api

  