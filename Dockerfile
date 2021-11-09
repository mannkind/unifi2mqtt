FROM mcr.microsoft.com/dotnet/sdk:6.0 as build
WORKDIR /src
COPY . .
RUN if [ ! -d output ]; then dotnet build -o output -c Release Unifi; fi

FROM mcr.microsoft.com/dotnet/runtime:6.0 AS runtime
COPY --from=build /src/output app
ENTRYPOINT ["dotnet", "./app/Unifi.dll"]
