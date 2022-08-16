# $BUILDPLATFORM ensures the native build platform is utilized
ARG BUILDPLATFORM=linux/amd64
FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:6.0 as build
WORKDIR /src
# Only fetch dependencies once
# Find the non-test csproj file, move it to the appropriate folder, and restore project deps
COPY Unifi/*.csproj ./Unifi/
RUN mkdir -p vendor && dotnet restore Unifi
COPY . ./
# Build the app
# Find the non-test csproj file, build that project
ARG BUILD_VERSION=0.0.0.0
RUN dotnet build -o output -c Release --no-restore -p:Version=$BUILD_VERSION Unifi

FROM mcr.microsoft.com/dotnet/runtime:6.0 AS runtime
COPY --from=build /src/output app
ENTRYPOINT ["dotnet", "./app/Unifi.dll"]
