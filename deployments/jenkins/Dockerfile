FROM jenkins/jenkins:2.516.1-lts
USER root
RUN apt update && apt install -y lsb-release ca-certificates curl && \
    install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc && \
    chmod a+r /etc/apt/keyrings/docker.asc && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
    https://download.docker.com/linux/debian $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable" \
    | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt update && apt-get install -y docker-ce-cli && \
    apt clean && rm -rf /var/lib/apt/lists/*
USER jenkins
# list jenkins plugins https://plugins.jenkins.io/
RUN jenkins-plugin-cli --plugins "blueocean docker-workflow json-path-api gitlab-plugin"
