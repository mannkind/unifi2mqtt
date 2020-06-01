FROM mcr.microsoft.com/dotnet/core/sdk:3.1 as build
WORKDIR /src
COPY . .
RUN if [ ! -d output ]; then dotnet build -o output -c Release Unifi; fi

FROM mcr.microsoft.com/dotnet/core/runtime:3.1 AS runtime
COPY --from=build /src/output app
ENTRYPOINT ["dotnet", "./app/Unifi.dll"]
