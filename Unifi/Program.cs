using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using KoenZomers.UniFi.Api;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT.Core;
using TwoMQTT.Core.DataAccess;
using TwoMQTT.Core.Extensions;
using Unifi.DataAccess;
using Unifi.Managers;
using Unifi.Models.Shared;


namespace Unifi
{
    class Program : ConsoleProgram<Resource, Command, SourceManager, SinkManager>
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            p.MapOldEnvVariables();
            await p.ExecuteAsync(args);
        }

        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            services.AddHttpClient<ISourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>>();

            return services
                .AddMemoryCache()
                .ConfigureOpts<Models.Shared.Opts>(hostContext, Models.Shared.Opts.Section)
                .ConfigureOpts<Models.SourceManager.Opts>(hostContext, Models.SourceManager.Opts.Section)
                .ConfigureOpts<Models.SinkManager.Opts>(hostContext, Models.SinkManager.Opts.Section)
                .AddTransient<Api>(x =>
                {
                    var opts = x.GetService<IOptions<Models.SourceManager.Opts>>();
                    return new Api(new Uri(opts.Value.Host));
                })
                .AddTransient<ISourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>>(x =>
                {
                    var opts = x.GetService<IOptions<Models.SourceManager.Opts>>();
                    return new SourceDAO(
                        x.GetService<ILogger<SourceDAO>>(), x.GetService<IMemoryCache>(), x.GetService<Api>(),
                        opts.Value.Username, opts.Value.Password, opts.Value.AwayTimeout
                    );
                });
        }

        [Obsolete("Remove in the near future.")]
        private void MapOldEnvVariables()
        {
            var found = false;
            var foundOld = new List<string>();
            var mappings = new[]
            {
                new { Src = "UNIFI_HOST", Dst = "UNIFI__HOST", CanMap = true, Strip = "", Sep = "" },
                new { Src = "UNIFI_USERNAME", Dst = "UNIFI__USERNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "UNIFI_PASSWORD", Dst = "UNIFI__PASSWORD", CanMap = true, Strip = "", Sep = "" },
                new { Src = "UNIFI_DEVICEMAPPING", Dst = "UNIFI__RESOURCES", CanMap = true, Strip = "",  Sep = ";" },
                new { Src = "UNIFI_LOOKUPINTERVAL", Dst = "UNIFI__POLLINGINTERVAL", CanMap = false, Strip = "", Sep = "" },
                new { Src = "UNIFI_AWAYTIMEOUT", Dst = "UNIFI__AWAYTIMEOUT", CanMap = false, Strip = "", Sep = "" },
                new { Src = "MQTT_TOPICPREFIX", Dst = "UNIFI__MQTT__TOPICPREFIX", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERY", Dst = "UNIFI__MQTT__DISCOVERYENABLED", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERYPREFIX", Dst = "UNIFI__MQTT__DISCOVERYPREFIX", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_DISCOVERYNAME", Dst = "UNIFI__MQTT__DISCOVERYNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_BROKER", Dst = "UNIFI__MQTT__BROKER", CanMap = true, Strip = "tcp://", Sep = "" },
                new { Src = "MQTT_USERNAME", Dst = "UNIFI__MQTT__USERNAME", CanMap = true, Strip = "", Sep = "" },
                new { Src = "MQTT_PASSWORD", Dst = "UNIFI__MQTT__PASSWORD", CanMap = true, Strip = "", Sep = "" },
            };

            foreach (var mapping in mappings)
            {
                var old = Environment.GetEnvironmentVariable(mapping.Src);
                if (string.IsNullOrEmpty(old))
                {
                    continue;
                }

                found = true;
                foundOld.Add($"{mapping.Src} => {mapping.Dst}");

                if (!mapping.CanMap)
                {
                    continue;
                }

                // Strip junk where possible
                if (!string.IsNullOrEmpty(mapping.Strip))
                {
                    old = old.Replace(mapping.Strip, string.Empty);
                }

                // Simple
                if (string.IsNullOrEmpty(mapping.Sep))
                {
                    Environment.SetEnvironmentVariable(mapping.Dst, old);
                }
                // Complex
                else
                {
                    var resourceSlugs = old.Split(",");
                    var i = 0;
                    foreach (var resourceSlug in resourceSlugs)
                    {
                        var parts = resourceSlug.Split(mapping.Sep);
                        var id = parts.Length >= 1 ? parts[0] : string.Empty;
                        var slug = parts.Length >= 2 ? parts[1] : string.Empty;
                        var idEnv = $"{mapping.Dst}__{i}__MACAddress";
                        var slugEnv = $"{mapping.Dst}__{i}__Slug";
                        Environment.SetEnvironmentVariable(idEnv, id);
                        Environment.SetEnvironmentVariable(slugEnv, slug);
                    }
                }

            }


            if (found)
            {
                var loggerFactory = LoggerFactory.Create(builder => { builder.AddConsole(); });
                var logger = loggerFactory.CreateLogger<Program>();
                logger.LogWarning("Found old environment variables.");
                logger.LogWarning($"Please migrate to the new environment variables: {(string.Join(", ", foundOld))}");
                Thread.Sleep(5000);
            }
        }
    }
}
