using TwoMQTT.Models;

namespace Unifi.Models.Options;

/// <summary>
/// The sink options
/// </summary>
public record MQTTOpts : MQTTManagerOptions
{
    public const string Section = "Unifi:MQTT";
    public const string TopicPrefixDefault = "home/unifi";
    public const string DiscoveryNameDefault = "unifi";
}
