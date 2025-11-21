# User Service

[![release](https://img.shields.io/github/v/release/qkitzero/user-service?logo=github)](https://github.com/qkitzero/user-service/releases)
[![test](https://github.com/qkitzero/user-service/actions/workflows/test.yml/badge.svg)](https://github.com/qkitzero/user-service/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/qkitzero/user-service/graph/badge.svg)](https://codecov.io/gh/qkitzero/user-service)
[![Buf CI](https://github.com/qkitzero/user-service/actions/workflows/buf-ci.yaml/badge.svg)](https://github.com/qkitzero/user-service/actions/workflows/buf-ci.yaml)

- Microservices Architecture
- gRPC
- gRPC Gateway
- Buf ([buf.build/qkitzero-org/user-service](https://buf.build/qkitzero-org/user-service))
- Clean Architecture
- Docker
- Test
- Codecov
- Cloud Build
- Cloud Run

```mermaid
flowchart TD
    subgraph gcp[GCP]
        secret_manager[Secret Manager]

        subgraph cloud_build[Cloud Build]
            build_user_service(Build user-service)
            push_user_service(Push user-service)
            deploy_user_service(Deploy user-service)

            build_user_service_gateway(Build user-service-gateway)
            push_user_service_gateway(Push user-service-gateway)
            deploy_user_service_gateway(Deploy user-service-gateway)
        end


        subgraph artifact_registry[Artifact Registry]
            user_service_image[(user-service image)]
            user_service_gateway_image[(user-service-gateway image)]
        end

        subgraph cloud_run[Cloud Run]
            user_service(User Service)
            user_service_gateway(User Service Gateway)
        end
    end

    subgraph external[External]
        auth_service(Auth Service)
        user_db[(User DB)]
    end

    build_user_service --> push_user_service --> user_service_image
    build_user_service_gateway --> push_user_service_gateway --> user_service_gateway_image

    user_service_image --> deploy_user_service --> user_service
    user_service_gateway_image --> deploy_user_service_gateway --> user_service_gateway

    secret_manager --> deploy_user_service
    secret_manager --> deploy_user_service_gateway

    user_service_gateway --> user_service
    user_service --> user_db
    user_service --> auth_service
```
