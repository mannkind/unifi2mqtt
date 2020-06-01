using TwoMQTT.Core.Models;

namespace Unifi.Models.SinkManager
{
    /// <summary>
    /// The sink options
    /// </summary>
    public class Opts : MQTTManagerOptions
    {
        public const string Section = "Unifi:MQTT";

        /// <summary>
        /// 
        /// </summary>
        public Opts()
        {
            this.TopicPrefix = "home/unifi";
            this.DiscoveryName = "unifi";
        }
    }
}
