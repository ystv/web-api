# Continuous Integration

We use Jenkins at YSTV to automate our build deployments. For each of our website 2020's repos we've included a `Jenkinsfile` which will allow you to replicate your own build environment like ours.

We use multi-stage pipelines and would recommend them for the best compatibility with our scripts.

The repo should be mostly plug and play with Jenkins, but you will need to set a couple of credentials that are used by all of our pipelines:

- `docker-registry-endpoint` (secret text) - A Docker registry endpoint i.e. `registry.ystv.co.uk`, check out this [guide](https://docs.docker.com/registry/) to learn more.

- `staging-server-address` (secret text) - Either an IP address or host name of a server running docker to SSH to
- `staging-server-path` (secret text) - The folder path of where all applications are stored i.e. `/opt`.
- `staging-server-key` (SSH Username with a private key) - A SSH key which will allow authentication to the staging server. Generate with `ssh-keygen`.

- `prod-server-address` (secret text) - Either an IP address or host name of a server running docker to SSH to.
- `prod-server-path` (secret text) - The folder path of where all applications are stored i.e. `/opt`.
- `prod-server-key` (SSH Username with a private key) - A SSH key which will allow authentication to the staging server. Generate with `ssh-keygen`.

Credentials specific to web-api

- `wapi-staging-env` (secret file) - A modified version of the `.env` file that would be used as your testing version.
- `wapi-prod-env` (secret file) - A modified version of the `.env` file that would be used as your public version.
