using Newtonsoft.Json;

namespace Unifi.Models.Source
{
    /// <summary>
    /// Hopefully this exists only until KoenZomers.UniFi.Api is updated.
    /// </summary>
    public record Clients
    {
        [JsonProperty(PropertyName = "mac")]
        public string MacAddress { get; set; } = string.Empty;
    }
}