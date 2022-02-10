# TODO: Minimize this image (900 --> 90 MB) by pulling
# FROM python:3-alpine
FROM python:3
# RUN pip install awscli
WORKDIR /
COPY main main
COPY conf/app.conf conf/app.conf
COPY entrypoint.sh entrypoint.sh
RUN chmod +x main entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
