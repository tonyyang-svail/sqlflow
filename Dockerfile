FROM sqlflow/sqlflow:latest

ARG NB_USER=jovyan
ARG NB_UID=1000
ENV USER ${NB_USER}
ENV NB_UID ${NB_UID}
ENV HOME /home/${NB_USER}

RUN adduser --disabled-password \
    --gecos "Default user" \
    --uid ${NB_UID} \
    ${NB_USER}

# Make sure the contents of our repo are in ${HOME}
COPY . ${HOME}
RUN chown -R ${NB_UID} ${HOME}

RUN echo "set -e" >> /entrypoint.sh
RUN echo "echo begin entrypoint.sh" >> /entrypoint.sh
RUN echo "service mysql start" >> /entrypoint.sh
RUN echo "su - jovyan" >> /entrypoint.sh
RUN echo "echo end entrypoint.sh" >> /entrypoint.sh
RUN echo "exec \"$@\"" >> /entrypoint.sh

RUN chmod +x /entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]

