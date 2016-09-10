FROM  ubuntu:trusty

# install slapd in noninteractive mode
RUN apt-get update && \
	echo 'slapd/root_password password password' | debconf-set-selections &&\
    echo 'slapd/root_password_again password password' | debconf-set-selections && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y slapd ldap-utils &&\
	rm -rf /var/lib/apt/lists/*

ADD files /ldap

RUN service slapd start ;\
    cd /ldap &&\
	ldapadd -Y EXTERNAL -H ldapi:/// -f back.ldif &&\
	ldapadd -Y EXTERNAL -H ldapi:/// -f sssvlv_load.ldif &&\
    ldapadd -Y EXTERNAL -H ldapi:/// -f sssvlv_config.ldif &&\
    ldapadd -x -D cn=admin,dc=pivotal,dc=org -w password -c -f front.ldif &&\
    ldapadd -x -D cn=admin,dc=pivotal,dc=org -w password -c -f more.ldif

EXPOSE 389

CMD slapd -h 'ldap:/// ldapi:///' -g openldap -u openldap -F /etc/ldap/slapd.d -d stats
