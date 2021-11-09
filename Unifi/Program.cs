using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using TwoMQTT;
using TwoMQTT.Extensions;
using TwoMQTT.Interfaces;
using TwoMQTT.Managers;
using Unifi.DataAccess;
using Unifi.Liasons;
using Unifi.Models.Options;
using Unifi.Models.Shared;

await ConsoleProgram<Resource, object, SourceLiason, MQTTLiason>.
    ExecuteAsync(args,
        envs: new Dictionary<string, string>()
        {
            {
                $"{MQTTOpts.Section}:{nameof(MQTTOpts.TopicPrefix)}",
                MQTTOpts.TopicPrefixDefault
            },
            {
                $"{MQTTOpts.Section}:{nameof(MQTTOpts.DiscoveryName)}",
                MQTTOpts.DiscoveryNameDefault
            },
        },
        configureServices: (HostBuilderContext context, IServiceCollection services) =>
        {
            services
                .AddMemoryCache()
                .AddOptions<SharedOpts>(SharedOpts.Section, context.Configuration)
                .AddOptions<SourceOpts>(SourceOpts.Section, context.Configuration)
                .AddOptions<TwoMQTT.Models.MQTTManagerOptions>(MQTTOpts.Section, context.Configuration)
                .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                {
                    var opts = x.GetRequiredService<IOptions<SourceOpts>>();
                    return new ThrottleManager(opts.Value.PollingInterval);
                })
                .AddTypeNamedHttpClient<ApiControllerDetection>(allowAutoRedirect: false)
                .AddTypeNamedHttpClient<Api>(lifetime: System.Threading.Timeout.InfiniteTimeSpan)
                .AddSingleton<Api>(x =>
                {
                    var opts = x.GetRequiredService<IOptions<SourceOpts>>();
                    var hcf = x.GetRequiredService<IHttpClientFactory>(); // Hopefully this only exists until KoenZomers.UniFi.Api is updated.
                    return new Api(new Uri(opts.Value.Host), opts.Value.Site, hcf);
                })
                .AddSingleton<ISourceDAO>(x =>
                {
                    var logger = x.GetRequiredService<ILogger<SourceDAO>>();
                    var cache = x.GetRequiredService<IMemoryCache>();
                    var api = x.GetRequiredService<Unifi.DataAccess.Api>();
                    var opts = x.GetRequiredService<IOptions<SourceOpts>>();
                    return new SourceDAO(logger,
                        cache,
                        api,
                        opts.Value.Username,
                        opts.Value.Password,
                        opts.Value.AwayTimeout);
                });
        });