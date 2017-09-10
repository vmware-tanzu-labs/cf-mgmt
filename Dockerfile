FROM concourse/buildroot:git

MAINTAINER Caleb Washburn "cwashburn@pivotal.io"

COPY cf-mgmt /usr/bin/cf-mgmt
RUN chmod +x /usr/bin/cf-mgmt && cf-mgmt --version
