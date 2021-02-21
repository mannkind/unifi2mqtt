using System.Collections.Generic;
using Newtonsoft.Json;

namespace Unifi.Models.Source
{
    /// <summary>
    /// Hopefully this exists only until KoenZomers.UniFi.Api is updated.
    /// </summary>
    public record Payload<T>
    {
        [JsonProperty(PropertyName = "data")]
        public List<T> Data { get; init; } = new();
    }
}