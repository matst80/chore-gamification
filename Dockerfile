# syntax=docker/dockerfile:1

FROM golang AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
#COPY pkg pkg
RUN CGO_ENABLED=1 GOOS=linux go build -o /chore-gamification

FROM gcr.io/distroless/base-debian11 
WORKDIR /

#EXPOSE 25
EXPOSE 8081

#COPY *.html /
COPY --from=build-stage /chore-gamification /chore-gamification
ENTRYPOINT ["/chore-gamification"]