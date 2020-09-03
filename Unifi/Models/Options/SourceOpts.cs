using System;

namespace Unifi.Models.Options
{
    /// <summary>
    /// The source options
    /// </summary>
    public record SourceOpts
    {
        public const string Section = "Unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Host { get; init; } = "https://unifi.local:8443";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Username { get; init; } = "unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Password { get; init; } = "unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Site { get; init; } = "default";

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan AwayTimeout { get; init; } = new TimeSpan(0, 5, 1);

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; init; } = new TimeSpan(0, 0, 11);
    }
}
