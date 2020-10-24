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
    class Program : ConsoleProgram<Resource, object, SourceLiason, MQTTLiason>
    {
        static async Task Main(string[] args)
        {
            var p = new Program();
            await p.ExecuteAsync(args);
        }

        /// <inheritdoc />
        protected override IDictionary<string, string> EnvironmentDefaults()
        {
            var sep = "__";
            var section = Models.Options.MQTTOpts.Section.Replace(":", sep);
            var sectsep = $"{section}{sep}";

            return new Dictionary<string, string>
            {
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.TopicPrefix)}", Models.Options.MQTTOpts.TopicPrefixDefault },
                { $"{sectsep}{nameof(Models.Options.MQTTOpts.DiscoveryName)}", Models.Options.MQTTOpts.DiscoveryNameDefault },
            };
        }

        /// <inheritdoc />
        protected override IServiceCollection ConfigureServices(HostBuilderContext hostContext, IServiceCollection services)
        {
            services.AddHttpClient<ISourceDAO>();

            return services
                .AddMemoryCache()
                .ConfigureOpts<Models.Options.SharedOpts>(hostContext, Models.Options.SharedOpts.Section)
                .ConfigureOpts<Models.Options.SourceOpts>(hostContext, Models.Options.SourceOpts.Section)
                .ConfigureOpts<TwoMQTT.Models.MQTTManagerOptions>(hostContext, Models.Options.MQTTOpts.Section)
                .AddSingleton<IThrottleManager, ThrottleManager>(x =>
                {
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();
                    if (opts == null)
                    {
                        throw new ArgumentException($"{nameof(opts.Value.PollingInterval)} is required for {nameof(ThrottleManager)}.");
                    }

                    return new ThrottleManager(opts.Value.PollingInterval);
                })
                .AddSingleton<Api>(x =>
                {
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();
                    if (opts == null)
                    {
                        throw new ArgumentException($"{nameof(opts.Value.Host)} is required for {nameof(Api)}.");
                    }

                    return new Api(new Uri(opts.Value.Host));
                })
                .AddSingleton<ISourceDAO>(x =>
                {
                    var logger = x.GetService<ILogger<SourceDAO>>();
                    var cache = x.GetService<IMemoryCache>();
                    var api = x.GetService<Api>();
                    var opts = x.GetService<IOptions<Models.Options.SourceOpts>>();

                    if (logger == null)
                    {
                        throw new ArgumentException($"{nameof(logger)} is required for {nameof(SourceDAO)}.");
                    }
                    if (cache == null)
                    {
                        throw new ArgumentException($"{nameof(cache)} is required for {nameof(SourceDAO)}.");
                    }
                    if (api == null)
                    {
                        throw new ArgumentException($"{nameof(api)} is required for {nameof(SourceDAO)}.");
                    }
                    if (logger == null)
                    {
                        throw new ArgumentException($"{nameof(logger)} is required for {nameof(SourceDAO)}.");
                    }
                    if (opts == null)
                    {
                        throw new ArgumentException($"{nameof(opts.Value.Username)}, {nameof(opts.Value.Password)}, and {nameof(opts.Value.AwayTimeout)} are required for {nameof(SourceDAO)}.");
                    }

                    return new SourceDAO(
                        logger, cache, api,
                        opts.Value.Username, opts.Value.Password, opts.Value.AwayTimeout
                    );
                });
        }
    }
}
