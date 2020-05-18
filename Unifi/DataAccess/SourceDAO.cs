using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Caching.Memory;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Primitives;
using KoenZomers.UniFi.Api;
using Newtonsoft.Json;
using TwoMQTT.Core.DataAccess;
using Unifi.Models.Shared;

namespace Unifi.DataAccess
{
    /// <summary>
    /// An class representing a managed way to interact with a source.
    /// </summary>
    public class SourceDAO : HTTPSourceDAO<SlugMapping, Command, Models.SourceManager.FetchResponse, object>
    {
        /// <summary>
        /// Initializes a new instance of the SourceDAO class.
        /// </summary>
        /// <param name="logger"></param>
        /// <param name="httpClientFactory"></param>
        /// <param name="cache"></param>
        /// <param name="unifiClient"></param>
        /// <param name="username"></param>
        /// <param name="password"></param>
        /// <param name="awayTimeout"></param>
        /// <returns></returns>
        public SourceDAO(ILogger<SourceDAO> logger, IHttpClientFactory httpClientFactory, IMemoryCache cache, Api unifiClient, string username, string password, TimeSpan awayTimeout) :
            base(logger, httpClientFactory)
        {
            this.Cache = cache;
            this.Username = username;
            this.Password = password;
            this.AwayTimeout = awayTimeout;
            this.UnifiClient = unifiClient;
            this.UnifiClient.DisableSslValidation();
        }

        /// <inheritdoc />
        public override async Task<Models.SourceManager.FetchResponse?> FetchOneAsync(SlugMapping data,
            CancellationToken cancellationToken = default)
        {
            try
            {
                return await this.FetchAsync(data.MACAddress, cancellationToken);
            }
            catch (Exception e)
            {
                var msg = e is HttpRequestException ? "Unable to fetch from the Unifi API" :
                          e is JsonException ? "Unable to deserialize response from the Unifi API" :
                          "Unable to send to the Unifi API";
                this.Logger.LogError(msg, e);
                this.IsLoggedIn = false;
                return null;
            }
        }

        /// <summary>
        /// The internal cache.
        /// </summary>
        private readonly IMemoryCache Cache;

        /// <summary>
        /// The Username to access the source.
        /// </summary>
        private readonly string Username;

        /// <summary>
        /// The Password to access the source.
        /// </summary>
        private readonly string Password;

        /// <summary>
        /// The client to access the source.
        /// </summary>
        private readonly Api UnifiClient;

        /// <summary>
        /// An internal timeout for how long until a device is considered away.
        /// </summary>
        private readonly TimeSpan AwayTimeout;

        /// <summary>
        /// The semaphore to limit how many times the source api is called.
        /// </summary>
        private readonly SemaphoreSlim ClientsSemaphore = new SemaphoreSlim(1, 1);

        /// <summary>
        /// A flag that indicates if logged into the source.
        /// </summary>
        private bool IsLoggedIn = false;

        /// <summary>
        /// 
        /// </summary>
        /// <typeparam name="string"></typeparam>
        /// <typeparam name="DateTime"></typeparam>
        /// <returns></returns>
        private readonly ConcurrentDictionary<string, DateTime> LastSeen = new ConcurrentDictionary<string, DateTime>();

        /// <summary>
        /// Fetch one response from the source
        /// </summary>
        /// <param name="data"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        private async Task<Models.SourceManager.FetchResponse?> FetchAsync(string macAddress,
            CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug($"Started finding {macAddress} from Unifi");
            var clients = await this.AllClientsAsync(cancellationToken);
            if (clients == null)
            {
                this.Logger.LogDebug($"Unable to find {macAddress} from Unifi");
                return null;
            }

            var client = clients.FirstOrDefault(x => x.MacAddress == macAddress);
            if (client != null)
            {
                this.LastSeen[client.MacAddress] = DateTime.Now;
            }

            var dt = DateTime.MinValue;
            if (this.LastSeen.ContainsKey(macAddress))
            {
                dt = this.LastSeen[macAddress];
            }

            return new Models.SourceManager.FetchResponse
            {
                MACAddress = macAddress,
                State = dt > (DateTime.Now - this.AwayTimeout),
            };
        }

        /// <summary>
        /// Fetch one response from the source
        /// </summary>
        /// <param name="data"></param>
        /// <param name="cancellationToken"></param>
        /// <returns></returns>
        private async Task<IEnumerable<KoenZomers.UniFi.Api.Responses.Clients>> AllClientsAsync(CancellationToken cancellationToken = default)
        {
            this.Logger.LogDebug($"Started fetching all clients from Unifi");
            await this.ClientsSemaphore.WaitAsync();

            try
            {
                // Check cache first to avoid hammering the API
                if (this.Cache.TryGetValue(ACTIVECLIENTS, out IEnumerable<KoenZomers.UniFi.Api.Responses.Clients> cachedObj))
                {
                    this.Logger.LogDebug($"Found all clients in the cache");
                    return cachedObj;
                }

                if (!this.IsLoggedIn)
                {
                    this.IsLoggedIn = await this.UnifiClient.Authenticate(this.Username, this.Password);
                }

                var clients = await this.UnifiClient.GetActiveClients();

                this.Logger.LogDebug($"Caching {clients.Count} clients");
                var cts = new CancellationTokenSource(new TimeSpan(0, 0, 9));
                var cacheOpts = new MemoryCacheEntryOptions()
                     .AddExpirationToken(new CancellationChangeToken(cts.Token));
                this.Cache.Set(ACTIVECLIENTS, clients, cacheOpts);
                return clients;
            }
            finally
            {
                this.ClientsSemaphore.Release();
            }
        }

        private const string ACTIVECLIENTS = "CLIENTS";
    }
}
