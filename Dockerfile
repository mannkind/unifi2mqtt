FROM mcr.microsoft.com/dotnet/core/sdk:3.1 as build
ARG TARGETPLATFORM
WORKDIR /src
COPY . .
RUN if [ ! -d output/$TARGETPLATFORM ]; then dotnet build -o output/$TARGETPLATFORM -c Release Unifi; fi \
    && cp -r output/$TARGETPLATFORM archoutput

FROM mcr.microsoft.com/dotnet/core/runtime:3.1 AS runtime
COPY --from=build /src/archoutput app
ENTRYPOINT ["./app/Unifi"]
