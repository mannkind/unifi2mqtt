using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using KoenZomers.UniFi.Api;
using TwoMQTT;
using TwoMQTT.Extensions;
using TwoMQTT.Interfaces;
using TwoMQTT.Managers;
using Unifi.DataAccess;
using Unifi.Liasons;
using Unifi.Models.Shared;


namespace Unifi
{
    class Program
    {
        static async Task Main(string[] args) => await ConsoleProgram<Resource, object, SourceLiason, MQTTLiason>.
            ExecuteAsync(args,
                envs: new Dictionary<string, string>()
                {
                    {
                        $"{Models.Options.MQTTOpts.Section}:{nameof(Models.Options.MQTTOpts.TopicPrefix)}",
                        Models.Options.MQTTOpts.TopicPrefixDefault
                    },
                    {
                        $"{Models.Options.MQTTOpts.Section}:{nameof(Models.Options.MQTTOpts.DiscoveryName)}",
                        Models.Options.MQTTOpts.DiscoveryNameDefault
                    },
                },
                configureServices: (HostBuilderContext context, IServiceCollection services) =>
                {
                    services
                        .AddMemoryCache()
                        .AddOptions<Models.Options.SharedOpts>(Models.Options.SharedOpts.Section, context.Configuration)
                        .AddOptions<Models.Options.SourceOpts>(Models.Options.SourceOpts.Section, context.Configuration)
                        .AddOptions<TwoMQTT.Models.MQTTManagerOptions>(Models.Options.MQTTOpts.Section, context.Configuration)
                        .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                        {
                            var opts = x.GetRequiredService<IOptions<Models.Options.SourceOpts>>();
                            return new ThrottleManager(opts.Value.PollingInterval);
                        })
                        .AddSingleton<Api>(x =>
                        {
                            var opts = x.GetRequiredService<IOptions<Models.Options.SourceOpts>>();
                            return new Api(new Uri(opts.Value.Host));
                        })
                        .AddSingleton<ISourceDAO>(x =>
                        {
                            var logger = x.GetRequiredService<ILogger<SourceDAO>>();
                            var cache = x.GetRequiredService<IMemoryCache>();
                            var api = x.GetRequiredService<Api>();
                            var opts = x.GetRequiredService<IOptions<Models.Options.SourceOpts>>();
                            return new SourceDAO(logger,
                                cache,
                                api,
                                opts.Value.Username,
                                opts.Value.Password,
                                opts.Value.AwayTimeout);
                        });
                });
    }
}