FROM 525143545807.dkr.ecr.eu-central-1.amazonaws.com/ubuntu-go-build:22.04 as builder

LABEL stage=gobuilder


# 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
#    GOPRIVATE=dghire.com \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

#

WORKDIR /application

# RUN git config --global url."https://$GIT_USER:$GIT_PWD@git-codecommit.eu-central-1.amazonaws.com/".insteadOf "https://git-codecommit.eu-central-1.amazonaws.com/"

#COPY go.mod , go.sum and download the dependencied
COPY . .
RUN go mod download

RUN go build -ldflags "-s -w" -o /application/build/html2pdf

#FROM public.ecr.aws/ubuntu/ubuntu:22.04_stable
FROM 525143545807.dkr.ecr.eu-central-1.amazonaws.com/ubuntu-go-run:22.04

WORKDIR /target

# 复制编译后的程序
COPY --from=builder /application/build/html2pdf /target/html2pdf
#COPY --from=builder /application/fonts/ /usr/share/fonts/
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
ENTRYPOINT ["/target/html2pdf"]
