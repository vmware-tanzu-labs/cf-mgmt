FROM concourse/buildroot:git

MAINTAINER Caleb Washburn "cwashburn@pivotal.io"

COPY cf-mgmt-linux /usr/bin/cf-mgmt
RUN chmod +x /usr/bin/cf-mgmt
RUN cf-mgmt --version
