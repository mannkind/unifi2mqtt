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
using Newtonsoft.Json;
using TwoMQTT.Interfaces;
using Unifi.Models.Shared;
using Unifi.Models.Source;

namespace Unifi.DataAccess;

public interface ISourceDAO : IPollingSourceDAO<SlugMapping, Response, object, object>
{
}

/// <summary>
/// An class representing a managed way to interact with a source.
/// </summary>
public class SourceDAO : ISourceDAO
{
    /// <summary>
    /// Initializes a new instance of the SourceDAO class.
    /// </summary>
    /// <param name="logger"></param>
    /// <param name="cache"></param>
    /// <param name="unifiClient"></param>
    /// <param name="username"></param>
    /// <param name="password"></param>
    /// <param name="awayTimeout"></param>
    /// <returns></returns>
    public SourceDAO(ILogger<SourceDAO> logger, IMemoryCache cache, Api unifiClient, string username, string password, TimeSpan awayTimeout)
    {
        this.Logger = logger;
        this.Cache = cache;
        this.Username = username;
        this.Password = password;
        this.AwayTimeout = awayTimeout;
        this.UnifiClient = unifiClient;
    }

    /// <inheritdoc />
    public async Task<Response?> FetchOneAsync(SlugMapping data,
        CancellationToken cancellationToken = default)
    {
        try
        {
            return await this.FetchAsync(data.MACAddress, cancellationToken);
        }
        catch (Exception e)
        {
            var msg = e switch
            {
                HttpRequestException => "Unable to fetch from the Unifi API",
                JsonException => "Unable to deserialize response from the Unifi API",
                _ => "Unable to send to the Unifi API"
            };
            this.Logger.LogError(msg + "; {exception}", e);
            this.IsLoggedIn = false;
            return null;
        }
    }

    /// <summary>
    /// The logger used internally.
    /// </summary>
    private readonly ILogger<SourceDAO> Logger;

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
    /// 
    /// </summary>
    private const string ACTIVECLIENTS = "CLIENTS";

    /// <summary>
    /// Fetch one response from the source
    /// </summary>
    /// <param name="data"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    private async Task<Response?> FetchAsync(string macAddress,
        CancellationToken cancellationToken = default)
    {
        this.Logger.LogDebug("Started finding {macAddress} from Unifi", macAddress);
        var clients = await this.AllClientsAsync(cancellationToken);
        if (clients == null)
        {
            this.Logger.LogDebug("Unable to find {macAddress} from Unifi", macAddress);
            return null;
        }

        var client = clients.FirstOrDefault(x => x.MacAddress.Equals(macAddress, StringComparison.OrdinalIgnoreCase));
        if (client != null)
        {
            this.LastSeen[macAddress] = DateTime.Now;
        }

        var dt = DateTime.MinValue;
        if (this.LastSeen.ContainsKey(macAddress))
        {
            dt = this.LastSeen[macAddress];
        }

        this.Logger.LogDebug("{macAddress} found in controller: {found}; last seen at {dt}; presence: {state}", macAddress, client != null, dt, dt > (DateTime.Now - this.AwayTimeout));
        return new Response
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
    private async Task<IEnumerable<Clients>> AllClientsAsync(CancellationToken cancellationToken = default)
    {
        this.Logger.LogDebug("Started fetching all clients from Unifi");
        await this.ClientsSemaphore.WaitAsync();

        try
        {
            // Check cache first to avoid hammering the API
            if (this.Cache.TryGetValue(ACTIVECLIENTS, out IEnumerable<Clients> cachedObj))
            {
                this.Logger.LogDebug("Found all clients in the cache");
                return cachedObj;
            }

            if (!this.IsLoggedIn)
            {
                this.IsLoggedIn = await this.UnifiClient.Authenticate(this.Username, this.Password);
            }

            var clients = await this.UnifiClient.GetActiveClients();

            this.Logger.LogDebug("Caching {count} clients", clients.Count);
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
}
