FROM concourse/buildroot:git

MAINTAINER Caleb Washburn "cwashburn@pivotal.io"

COPY cf-mgmt-linux /usr/bin/cf-mgmt
COPY cf-mgmt-config-linux /usr/bin/cf-mgmt-config
RUN chmod +x /usr/bin/cf-mgmt && chmod +x /usr/bin/cf-mgmt-config
RUN cf-mgmt version
