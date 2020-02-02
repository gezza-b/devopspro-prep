# Dockerfile for Jenkins Master FROM Jenkins
# Add the entry amazon-ecs to plugin.txt to preload the Amazon ECS plugin
COPY plugins.txt /usr/share/jenkins/plugins.txt
RUN /usr/local/bin/plugins.sh /usr/share/jenkins/plugins.txt