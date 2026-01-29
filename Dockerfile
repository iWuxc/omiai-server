FROM isheji-registry-vpc.cn-beijing.cr.aliyuncs.com/huaxia/golang:1.25.4

ENV APP_HOME=/app
ENV TZ=Asia/Shanghai

WORKDIR $APP_HOME

COPY . /app

RUN  make build

RUN  mv /app/configs/config.yaml.bak /app/configs/config.yaml
RUN  mv /app/configs/nacos.yaml.bak /app/configs/nacos.yaml

CMD ["sh", "-c", "/app/bin/server -conf=/app/configs"]