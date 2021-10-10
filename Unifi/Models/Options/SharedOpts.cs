using System.Collections.Generic;
using TwoMQTT.Interfaces;
using Unifi.Models.Shared;

namespace Unifi.Models.Options
{
    /// <summary>
    /// The shared options across the application
    /// </summary>
    public record SharedOpts : ISharedOpts<SlugMapping>
    {
        public const string Section = "Unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <typeparam name="SlugMapping"></typeparam>
        /// <returns></returns>
        public List<SlugMapping> Resources { get; init; } = new();

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool AsDeviceTracker { get; init; } = false;
    }
}
