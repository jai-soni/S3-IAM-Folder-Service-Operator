FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/team2-kubeop \
    USER_UID=1001 \
    USER_NAME=team2-kubeop

# install operator binary
COPY build/_output/bin/team2-kubeop ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
