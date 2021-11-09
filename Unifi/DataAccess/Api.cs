using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Newtonsoft.Json;
using Unifi.Models.Source;

namespace Unifi.DataAccess;

/// <summary>
/// Hopefully this exists only until KoenZomers.UniFi.Api is updated.
/// </summary>
public class Api
{
    public Api(Uri host, string site, IHttpClientFactory httpClientFactory)
    {
        this.Host = host;
        this.Site = site;
        this.HttpClientFactory = httpClientFactory;
    }

    public async Task<bool> Authenticate(string username, string password, CancellationToken cancellationToken = default)
    {
        using var dclient = this.HttpClientFactory.CreateClient(nameof(ApiControllerDetection));
        var detectResp = await dclient.GetAsync(this.Host, cancellationToken);
        this.UniFiOS = detectResp.StatusCode == System.Net.HttpStatusCode.OK;

        using var data = new StringContent(JsonConvert.SerializeObject(new
        {
            username = username,
            password = password,
            remember = false,
        }), null, "application/json");

        using var client = this.HttpClientFactory.CreateClient(nameof(Api));
        var resp = await client.PostAsync($"{this.Host}{this.MapUrl($"api/login")}", data, cancellationToken);

        return resp.IsSuccessStatusCode;
    }

    public async Task<List<Unifi.Models.Source.Clients>> GetActiveClients(CancellationToken cancellation = default)
    {
        using var client = this.HttpClientFactory.CreateClient(nameof(Api));
        var resp = await client.GetAsync($"{this.Host}{this.MapUrl($"api/s/{this.Site}/stat/sta")}", cancellation);
        var result = await resp.Content.ReadAsStringAsync(cancellation);
        var objs = JsonConvert.DeserializeObject<Payload<Clients>>(result);

        return objs.Data;
    }

    /// <summary>
    /// 
    /// </summary>
    private readonly Uri Host;

    /// <summary>
    /// 
    /// </summary>
    private readonly string Site;

    /// <summary>
    /// 
    /// </summary>
    private readonly IHttpClientFactory HttpClientFactory;

    /// <summary>
    /// 
    /// </summary>
    private bool UniFiOS;

    private string MapUrl(string url)
    {
        if (!this.UniFiOS)
        {
            return url;
        }

        if (url == "api/login")
        {
            return "api/auth/login";
        }

        return $"proxy/network/{url}";
    }
}

/// <summary>
/// Hopefully this exists only until KoenZomers.UniFi.Api is updated.
/// </summary>
public class ApiControllerDetection { }
