FROM sqlflow/sqlflow:latest

ARG NB_USER=jovyan
ARG NB_UID=1000
ENV USER ${NB_USER}
ENV NB_UID ${NB_UID}
ENV HOME /home/${NB_USER}

RUN adduser --disabled-password \
    --gecos "Default user" \
    --uid 1000 jovyan

# autoload IPython magic command
ENV IPYTHON_STARTUP ${HOME}/.ipython/profile_default/startup/
RUN mkdir -p $IPYTHON_STARTUP
RUN echo 'get_ipython().magic(u"%reload_ext sqlflow.magic")' >> $IPYTHON_STARTUP/00-first.py
RUN echo 'get_ipython().magic(u"%autoreload 2")' >> $IPYTHON_STARTUP/00-first.py

# Make sure the contents of our repo are in ${HOME}
COPY . ${HOME}
RUN chown -R ${NB_UID} ${HOME}

# Activate Python virtual environment sqlflow-dev
# RUN echo "export PATH=/miniconda/envs/sqlflow-dev/bin:/miniconda/bin:$PATH" >> ${HOME}/.bashrc
# RUN echo "source /miniconda/bin/activate sqlflow-dev" >> ${HOME}/.bashrc

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["bash", "/entrypoint.sh"]

