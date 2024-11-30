# UALA_CHALLENGE

Challenge técnico para Dev Backend en Ualá

## Ejecución

La demo puede ser ejecutada de manera local de la siguiente forma:

Lanzar el servidor web:

```bash
docker compose -f "docker/docker-compose.yml" up -d
```

En este mismo proyecto se encuentra la collection de Postman lista para ser utilizada.

La misma se llama "Uala Challenge.postman_collection.json"

## Api Docs

La documentación de los endpoints disponibles se encuentra en el documento `swagger.json`.

Se puede ver de forma sencilla con Swagger UI. Para eso, se debe correr el servidor web:

```bash
docker run -p 8081:8080 -e SWAGGER_JSON=/tmp/swagger.yaml -v $PWD/docs/:/tmp swaggerapi/swagger-ui
```

Al ingresar a [localhost:8081](localhost:8081) se podrá visualizar.