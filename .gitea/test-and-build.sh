#!/bin/sh

curl \
--fail-with-body \
-H "Authorization: ${RUNNER_API_KEY}" \
-XPOST 192.168.1.186 \
-H 'Content-Type: application/json; charset=utf-8' \
--data-binary @- <<EOF
{
	"commands": [
		"git clone --single-branch --branch ${GITHUB_REF_NAME} --depth 1 ssh://git@gitea.larvit.se:21022/pwrpln/auth-api.git",
        "cd auth-api",
        "docker compose build",
        "docker compose --profile tests build",
        "docker compose run --rm tests",
        "docker compose down -v --remove-orphans -t0",
		"echo \"${DOCKER_PASSWORD}\" | docker login gitea.larvit.se --username ${DOCKER_USERNAME} --password-stdin",
		"docker build -t gitea.larvit.se/pwrpln/auth-api:${GITHUB_REF_NAME} .",
		"docker push gitea.larvit.se/pwrpln/auth-api:${GITHUB_REF_NAME}"
	]
}
EOF
