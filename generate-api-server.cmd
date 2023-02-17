go install golang.org/x/tools/cmd/goimports@latest

docker run --rm ^
    -v "%cd%":/local ^
    openapitools/openapi-generator-cli:v6.3.0 generate ^
    -g go-server ^
    --git-user-id eliona-smart-building-assistant ^
    --git-repo-id thingdust-app ^
    -i /local/openapi.yaml ^
    -o /local/apiserver ^
    --additional-properties="packageName=apiserver,sourceFolder=,outputAsLibrary=true"

docker run --rm ^
    -v "%cd%":/local ^
    openapitools/openapi-generator-cli:v6.3.0 generate ^
    -g openapi ^
    -i /local/openapi.yaml ^
    -o /local/apiserver ^
    --additional-properties=outputFile=openapi.json

goimports -w ./apiserver