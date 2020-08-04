using System;

namespace Unifi.Models.Options
{
    /// <summary>
    /// The source options
    /// </summary>
    public class SourceOpts
    {
        public const string Section = "Unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Host { get; set; } = "https://unifi.local:8443";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Username { get; set; } = "unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Password { get; set; } = "unifi";

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Site { get; set; } = "default";

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan AwayTimeout { get; set; } = new TimeSpan(0, 5, 1);

        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public TimeSpan PollingInterval { get; set; } = new TimeSpan(0, 0, 11);
    }
}
